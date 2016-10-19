package tag

import (
	"testing"

	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
	"github.com/franela/goblin"
)

var repo = types.Repo{Owner: "owner", Name: "repo"}
var build = types.Build{Number: 0, Workspace: ".", Commit: "01234567890", Branch: "my_Branch"}
var masterbuild = types.Build{Number: 0, Workspace: ".", Commit: "01234567890", Branch: "master"}

var nodeProjectDev = types.Plugin{
	Repo: repo, Build: build, ProjectType: "node",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

var dotnetProjectDev = types.Plugin{
	Repo: repo, Build: build, ProjectType: "dotnet-core",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

var otherProjectDev = types.Plugin{
	Repo: repo, Build: build, ProjectType: "Malbolge",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

var dotnetProjectMaster = types.Plugin{
	Repo: repo, Build: masterbuild, ProjectType: "dotnet-core",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

var nodeProjectMaster = types.Plugin{
	Repo: repo, Build: masterbuild, ProjectType: "node",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

var otherProjectMaster = types.Plugin{
	Repo: repo, Build: masterbuild, ProjectType: "Malbolge",
	DockerStorageDriver: "overlay",
	DockerHubRepo:       "repo",
	DockerHubUser:       "user",
	DockerHubPass:       "secret",
	DockerHubEmail:      "example@example.com",
	GithubAccessToken:   "supersecret",
	RancherCatalogRepo:  "catalog",
	RancherCatalogName:  "repo",
	DryRun:              false,
}

func TestHookImage(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Tag", func() {
		g.It("Check the proper tags are built", func() {
			if tags, err := CreateDockerImageTags(nodeProjectDev); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"owner_repo_my-Branch_1.0.0_0_0123456"})
			}
			if tags, err := CreateDockerImageTags(dotnetProjectDev); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"owner_repo_my-Branch_1.0.1_0_0123456"})
			}
			if tags, err := CreateDockerImageTags(otherProjectDev); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"owner_repo_my-Branch_0_0123456"})
			}
			if tags, err := CreateDockerImageTags(nodeProjectMaster); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"v1.0.0", "latest"})
			}
			if tags, err := CreateDockerImageTags(dotnetProjectMaster); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"v1.0.1", "latest"})
			}
			if tags, err := CreateDockerImageTags(otherProjectMaster); true {
				g.Assert(err).Equal(nil)
				g.Assert(tags).Equal([]string{"master_0_0123456", "latest"})
			}
		})
	})
}
