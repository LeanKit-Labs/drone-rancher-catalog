package tag

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
)

//vars for reading semVer data
type projectJSON struct {
	Version string `json:"version"`
}

var fileInfo struct {
	Version string `json:"version"`
}

func getJSONVersionReader(fname string) func() (string, error) {
	return func() (string, error) {
		fileData, err := os.Open(fname)
		if err != nil {
			return "", err
		}

		jsonParser := json.NewDecoder(fileData)
		if err = jsonParser.Decode(&fileInfo); err != nil {
			return "", err
		}

		return fileInfo.Version, nil
	}
}

func replaceUnderscores(str string) string {
	return strings.Replace(str, "_", "-", -1)
}

//CreateDockerImageTags takes plugin information and returns a list
//of tags to use when publishing the project image to Docker Hub
//TODO: this function might not need to take an error
func CreateDockerImageTags(p types.Plugin) ([]string, error) {

	var projectMap = map[string]func() (string, error){
		"node":        getJSONVersionReader(fmt.Sprintf("%s/package.json", p.Workspace)),
		"dotnet-core": getJSONVersionReader(fmt.Sprintf("%s/project.json", p.Workspace)),
	}

	//read version
	version := ""
	if getVersion, ok := projectMap[p.ProjectType]; ok {
		if val, err := getVersion(); err == nil {
			version = val
		} else {
			return []string{}, err
		}
	}

	//handle master tag
	if p.Branch == "master" {
		if version != "" { //default master format is v<smver> of the version is present or master_build_shaw if version is not found
			return []string{fmt.Sprintf("v%s", version), "latest"}, nil
		}
		return []string{fmt.Sprintf("master_%d_%s", p.Build.Number, p.Build.Commit[:7]), "latest"}, nil
	}

	//return the long tag
	//githubOwner_githubRepo_branch_semVer_globalProjectNum_shortCommitSHA
	if version != "" {
		return []string{fmt.Sprintf("%s_%s_%s_%s_%d_%s", p.Repo.Owner, p.Repo.Name, p.Build.Branch, version, p.Build.Number, p.Build.Commit[:7])}, nil
	}
	return []string{fmt.Sprintf("%s_%s_%s_%d_%s", p.Repo.Owner, p.Repo.Name, p.Build.Branch, p.Build.Number, p.Build.Commit[:7])}, nil
}
