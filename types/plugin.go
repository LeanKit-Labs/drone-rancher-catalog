package types

//Plugin contains data needed for the plugin to run
type Plugin struct {
	Repo
	Build

	ProjectType         string `json:"project_type"`
	DockerStorageDriver string `json:"docker_storage_driver"`
	DockerHubRepo       string `json:"docker_hub_repo"`
	DockerHubUser       string `json:"docker_hub_user"`
	DockerHubPass       string `json:"docker_hub_pass"`
	DockerHubEmail      string `json:"docker_hub_email"`
	GithubAccessToken   string `json:"github_access_token"`
	RancherCatalogRepo  string `json:"rancher_catalog_repo"`
	RancherCatalogName  string `json:"rancher_catalog_name"`
	DryRun              bool   `json:"dry_run"`
}
