package docker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
)

const daemonStoragePath = "/drone/docker"
const dockerCmd string = "/usr/bin/docker"
const dockerFilename = "Dockerfile"
const registry = "https://index.docker.io/v1/"
const dockerContext = "."

//ew global, the intent is that this is set by a single exported function (like a c_tor)
var workingDir = ""

//PublishImage builds a docker image and publishes it to docker hub
//TODO workspace could just be the Dockerfile path
func PublishImage(image string, imageTags []string, p types.Plugin) error {

	workingDir = p.Workspace
	fmt.Println("starting daemon")
	if err := startDaemon(p.DockerStorageDriver); err != nil {
		return err
	}

	fmt.Println("building image")
	for _, tag := range imageTags {

		fullImageName := fmt.Sprintf("%s:%s", image, tag)
		if err := buildImage(fullImageName); err != nil {
			return err
		}

		//push to docker hub, could be done asynchronously
		if !p.DryRun {
			fmt.Println("docker login")
			if err := login(p.DockerHubUser, p.DockerHubPass, p.DockerHubEmail); err != nil {
				return err
			}

			fmt.Println("pushing image to docker hub")
			if err := pushImage(image); err != nil {
				return err
			}
		}
	}

	return nil
}

func startDaemon(storageDriver string) error {

	args := []string{"daemon", "-g", daemonStoragePath}

	if storageDriver != "" {
		args = append(args, "-s", storageDriver)
	}

	cmd := createCmd(args, false)

	//start the daemon in the background
	go func() {
		cmd.Run() //this never returns :(
	}()

	//poll until daemon is available or throw error
	isUp := false
	for i := 1; i <= 90; i++ {
		if err := createCmd([]string{"info"}, true).Run(); err == nil {
			isUp = true
			time.Sleep(1 * time.Second)
			break
		}
	}

	if !isUp {
		createCmd([]string{"info"}, false).Run()
		return errors.New("Timeout exceeded while starting docker daemon")
	}

	return nil

}

func login(dockerUser string, dockerPass string, dockerEmail string) error {
	args := []string{
		"login",
		"-u", dockerUser,
		"-p", dockerPass,
		"-e", dockerEmail, registry,
	}
	return createCmd(args, false).Run()
}

func buildImage(image string) error {
	args := []string{
		"build",
		"--pull=true",
		"--rm=true",
		"-f", dockerFilename,
		"-t", image,
		".",
	}

	return createCmd(args, true).Run()
}

func pushImage(image string) error {
	args := []string{
		"push",
		image,
	}
	return createCmd(args, false).Run()
}

//helper for executing shell commands
func createCmd(args []string, supressIO bool) *exec.Cmd {
	cmd := exec.Command(dockerCmd, args...)
	cmd.Dir = workingDir

	if supressIO {
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
