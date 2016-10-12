package tag

import (
	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
)

/*
	vars for reading semVer data
	type projectJSON struct {
		Version string `json:"version"`
	}
	var projectMap = map[string]string{
		"node":        "package.json",
		"dotnet-core": "project.json",
	}
*/

//CreateDockerImageTags takes plugin information and returns a list
//of tags to use when publishing the project image to Docker Hub
//TODO: this function might not need to take an error
func CreateDockerImageTags(p types.Plugin) ([]string, error) {
	/*TODOS:
				  1) if the branch being built is master 2 tags should be returned
		        -> latest
		        -> v<semVer>

		        if the semver can not be determined, then just return latest

		      2) if the branch being built is not master then return 1 tag of the form
		         -> githubOwner_githubRepo_branch_semVer_globalProjectNum_shortCommitSHA

	           shortCommitSHA can be just a 7 char sub string (git rev-parse --short)

		         if the semver can not be determined then it should not be included in the tag
	*/

	return []string{}, nil
}
