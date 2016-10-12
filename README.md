# drone-rancher-catalog
## __This plugin is in the protoype phase and not suitable for production__

### Goals

The idea behing this plugin is to enable deployment of applications via Docker Hub and Rancher stacks.
At a high level, this plugin will do the following after a successful Drone CI build

* Generate tag(s) for the resulting image
* Build Docker images for each tag
* Push the images to Docker Hub
* Generate a Rancher catalog entry from a template defined in it's repository

### Docker Image Builds

The plugin pushes 2 different kinds of images depending on the build context

#### Release images

If the branch being built is ```master``` the plugin build an image with the ```latest``` tag and one with a semver tag

EX:
* someRepo/someImage:latest
* someRepo/somImage:v2.0.0

#### Test Images

If the branch is not master, a single image is pushed with a tag in the following format

__someRepo/someImage:githubOwner\_githubRepo\_branch\_semVer\_globalBuildNumber\_shortCommitSHA__

The short commit SHA is a 7 char substring of the full commit ala ```git rev-parse HEAD --short```

Ex:

* leankit/myRepo:leankit_myProject_feature-x_1.0.0_22_absd3f5

### Detecting project version

If the plugin is given a __project_type__ argument, it will attempt to look for a version in the following location

* Node.JS -> __package.json__ at the project root
* .NET Core -> __project.json__ at the project root

__Note__ support for other projects is planned to be added over time

If the project being built is unsupported, or if there is an issue obtaining the version the following will occurr

* Release Build (master) -> only the ```latest``` tag will be pushed
* Test Builds -> the version will be excluded from the tag entirely

### Requirements

* Go 1.6+ (for vendoring)
* Drone .4.0 (Environment variable plugin args via environment variables are not supported)
* (Glide package manager)[https://glide.sh/]

### How to build

* ```glide install``` -> install dependencies in /vendor
* ```go build```

### How to run (Example)

```
./drone-cowpoke-catalog-publish <<EOF
{
	"build": {
		"Number": 1,
		"Branch": "my-feature",
		"Commit": "commitSHA"
	},
	"workspace": {
		"path": "root-dir-of-taget-application"
	},
	"repo": {
		"name": "github-repo",
		"owner": "github-repo-owner"
	},
	"vargs": {
    "project_type": "node",
		"docker_hub_repo": "someUser/someImage",
		"docker_hub_user": "someUser",
		"docker_hub_pass": "somePass",
		"github_access_token": "GitHub Outh Token",
		"rancher_catalog_repo": "someOwner/someRepo",
		"rancher_catalog_name": "name of catalog in Rancher",
		"dry_run": true
	}
}
EOF
```

### Plugin Options

|  Name | What It Does | Required?
|---|---|---|
| project_type | the type of project being built. Used to as a way to pull sem ver info from a project   | If ommitted, no semver information will be included in docker image tags
| docker_hub_repo | repo in docker hub where images will be pushed  | Y
| docker_hub_user  | username for pushing docker images  | Y
| docker_hub_pass | password of docker hub user| Y
| github_access_token | github access token used to push catalog entries, should have read/write perms on the catalog repo| Y
| rancher_catalog_repo | the github repository of the Rancher catatlog | Y
| rancher_catalog_name | name of the catalog in Rancher | Y
| dry_run | will simulate steps, but not push anything to docker hub or rancher. userful for debugging | Defaults to ```false```


