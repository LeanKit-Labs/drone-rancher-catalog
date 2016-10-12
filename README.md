# drone-rancher-catalog
Drone plugin to publish Docker images and create Rancher catalog entries

### This plugin is in the protoype phase and not suitable for production

### Requirements

* Go 1.6+
* Drone .4.0
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

|  Name | What It Does  |   Required?
|---|---|---|
| project_type | the type of project being built. Used to as a way to pull sem ver info from a project   | No. If ommitted, no semver information will be included in docker image tags
| docker_hub_repo | repo in docker hub where images will be pushed  | Yes.
| docker_hub_user  | username for pushing docker images  |
| docker_hub_pass | |
| github_access_token | |
| rancher_catalog_repo | |
| rancher_catalog_name | |


