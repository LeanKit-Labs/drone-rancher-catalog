package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
	"github.com/asaskevich/govalidator"
	"github.com/google/go-github/github"
)

type CatalogInfo struct {
	catalogRepo string
	version     int
	branch      string
}

//Template a data struture to store the generic catalog template
type GenericTemplate struct {
	Config         string
	DockerCompose  string
	RancherCompose string
	Icon           []byte
}

//BuiltTemplate a data struture to store the catalog template with agruments
type BuiltTemplate struct {
	branch         string
	tag            string
	Config         string
	DockerCompose  string
	RancherCompose string
	Icon           []byte
}

type folder struct {
	Name  string      `json:"name"`
}

type tmplArguments struct {
	Branch string
	Tag    string
	Count  int
}

type templateFile struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

//StringWriter write to a string
type StringWriter struct {
	Value *string
}

func (w StringWriter) Write(p []byte) (n int, err error) {
	*w.Value = string(p)
	return len(*w.Value), nil
}

func getBytesFromURL(client *http.Client, url string, token string) ([]byte, int, error) {
	//build request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, -1, nil
	}
	request.SetBasicAuth(token, "x-oauth-basic")
	request.Close = true

	//run request
	response, err := client.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	//parse response
	contents, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	return contents, response.StatusCode, nil
}

//NewGenericTemplate gets the Template from github
func NewGenericTemplate(owner string, repo string, token string) (*GenericTemplate, error) {
	//build url
	templateFolderURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/base", owner, repo)
	if !govalidator.IsURL(templateFolderURL) {
		return nil, errors.New("Github Owner and Repo invalid")
	}

	//get directory
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	var templateDir []templateFile
	contents, _, err := getBytesFromURL(client, templateFolderURL, token)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(contents, &templateDir)

	//download files
	var result GenericTemplate
	for _, file := range templateDir {
		switch file.Name {
		case "catalogIcon.png":
			fileContents, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, _, err
			}
			result.Icon = fileContents
		case "config.tmpl":
			fileContents, _, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.Config = string(fileContents)
		case "docker-compose.tmpl":
			fileContents, _, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.DockerCompose = string(fileContents)
		case "rancher-compose.tmpl":
			fileContents, _, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.DockerCompose = string(fileContents)
		}
	}

	if len(result.Icon) == 0 {
		return nil, errors.New("catalogIcon.png not found")
	}
	if len(result.DockerCompose) == 0 {
		return nil, errors.New("docker-compose.tmpl not found")
	}
	if len(result.RancherCompose) == 0 {
		return nil, errors.New("rancher-compose.tmpl not found")
	}
	if len(result.Config) == 0 {
		return nil, errors.New("config.tmpl not found")
	}

	return &result, nil

}

func fixTemplate(p *types.Plugin, name string, templatedString, tag string) (string, error) {
	tmpl, err := template.New(name).Parse(templatedString)
	if err != nil {
		return "", err
	}

	var args tmplArguments
	args.Branch = p.Branch
	args.Count = p.Build.Number
	args.Tag = tag

	var writer StringWriter
	if err := tmpl.Execute(writer, p); err != nil {
		return "", err
	}
	return *writer.Value, nil
}

//SubBuildInfo replaces the templated values in the file
func (t *GenericTemplate) SubBuildInfo(p *types.Plugin, tag string) (*BuiltTemplate, error) {

	var final BuiltTemplate
	final.branch = p.Branch
	final.tag = tag

	val1, err1 := fixTemplate(p, "docker-compose.yml", t.DockerCompose, tag)
	if err1 != nil {
		return nil, err1
	}
	final.DockerCompose = val1
	val2, err2 := fixTemplate(p, "rancher-compose.yml", t.RancherCompose, tag)
	if err2 != nil {
		return nil, err2
	}
	final.RancherCompose = val2
	val3, err3 := fixTemplate(p, "config.yml", t.Config, tag)
	if err3 != nil {
		return nil, err3
	}
	final.Config = val3
	return final, nil
}

func getTemplateNum(client string, url string, token string) (int, err){
	folderBytes, code, err := getBytesFromURL(client, url, token)
	if err != {
		return -1, err
	}
	if (code == 404) {
		return 0, nil
	}
	var folders []folder
	currentTemplate := -1 //empty folder
	if json.Unmarshal(folderBytes, &folders) != nil {
		for _,folder := range folders {
			number := strconv.Atoi(folder.Name)
			if ( number > currentTemplate){
				currentTemplate = number
			}
		}
	}
	return currentTemplate+1, nil
	

}

func commitFile(githubClient github.Client, owner string, repo string, path string, contents []byte, message string) error {
	branch := "master"
	opts := github.RepositoryContentFileOptions{
		Message: &message,
		Content: contents,
		Branch:  &branch,
	}
	_, _, err := githubClient.Repositories.CreateFile(owner, repo, path, &opts)
	if err != nil {
		return err
	}
	return nil
}

//Commit commits the file to github
func (t *BuiltTemplate) Commit(token string, owner string, repo string, buildNum int)(CatalogInfo, error){
	token := oauth2.Token{AccessToken: token}
	tokenSource := oauth2.StaticTokenSource(&token)
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	githubClient := github.NewClient(oauthClient)

	client := &http.Client{
		Timeout: time.Second * 60,
	}
	number, err := getTemplateNum(client, fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/templates/%s", owner, repo, t.branch), token)
	if err != nil {
		return nil, err
	}
	if err = commitFile(githubClient, owner, repo, fmt.Sprintf("templates/%s/catalogIcon.png", t.branch), t.Icon, fmt.Sprintf("Drone Build #%d: Changine Icon", buildNum)); err != nil {
		return nil, err
	}
	if err = commitFile(githubClient, owner, repo, fmt.Sprintf("templates/%s/config.yml", t.branch), []byte(t.Config)], fmt.Sprintf("Drone Build #%d: Changine config.yml", buildNum)); err != nil {
		return nil, err
	}
	if err = commitFile(githubClient, owner, repo, fmt.Sprintf("templates/%s/%d/docker-compose.yml", t.branch), []byte(t.DockerCompose)], fmt.Sprintf("Drone Build #%d: Changine docker-compose.yml", buildNum)); err != nil {
		return nil, err
	}
	if err = commitFile(githubClient, owner, repo, fmt.Sprintf("templates/%s/%d/rancher-compose.yml", t.branch), []byte(t.RancherCompose)], fmt.Sprintf("Drone Build #%d: Changine rancher-compose.yml", buildNum)); err != nil {
		return nil, err
	}
	var info CatalogInfo;
	info.catalogRepo = fmt.Sprint("%s/%s", owner, repo)
	info.version = number
	info.branch = t.branch
	return info nil
}
