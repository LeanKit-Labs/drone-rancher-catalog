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

const daemonStoragePath = "/var/lib/docker"
const dockerCmd string = "/usr/local/bin/docker"
const dockerFilename = "Dockerfile"
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
			if err := login(p.DockerHubUser, p.DockerHubPass); err != nil {
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
	doneChan := make(chan error)
	args := []string{"daemon", "-g", daemonStoragePath}

	if storageDriver != "" {
		args = append(args, "-s", storageDriver)
	}

	cmd := createCmd(args, true)

	//start the daemon in the background
	go func() {
		err := cmd.Run()
		if err != nil {
			doneChan <- err
		}

		//poll until daemon is available or throw error
		isUp := false
		for i := 1; i <= 3; i++ {
			if err := createCmd([]string{"info"}, false).Run(); err == nil {
				isUp = true
				break
			}

			time.Sleep(time.Second * time.Duration(i))
		}

		if isUp {
			doneChan <- nil
		} else {
			doneChan <- errors.New("Timeout exceeded while starting docker daemon")
		}

		close(doneChan)
	}()

	err := <-doneChan

	return err
}

func login(dockerUser string, dockerPass string) error {
	args := []string{
		"login",
		"-u", dockerUser,
		"-p", dockerPass,
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

	return createCmd(args, false).Run()
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
