package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/LeanKit-Labs/drone-rancher-catalog/types"
	"github.com/asaskevich/govalidator"
)

//Template a data struture to store the catalog template
type Template struct {
	Config         string
	DockerCompose  string
	RancherCompose string
	Icon           []byte
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

func getBytesFromURL(client *http.Client, url string, token string) ([]byte, error) {
	//build request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
		return nil, nil
	}
	request.SetBasicAuth(token, "x-oauth-basic")
	request.Close = true

	//run request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	//parse response
	contents, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	return contents, nil
}

//NewTemplate gets the Template from github
func NewTemplate(owner string, repo string, token string) (*Template, error) {
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
	contents, err := getBytesFromURL(client, templateFolderURL, token)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(contents, &templateDir)

	//download files
	var result Template
	for _, file := range templateDir {
		switch file.Name {
		case "catalogIcon.png":
			fileContents, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.Icon = fileContents
		case "config.tmpl":
			fileContents, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.Config = string(fileContents)
		case "docker-compose.tmpl":
			fileContents, err := getBytesFromURL(client, file.DownloadURL, token)
			if err != nil {
				return nil, err
			}
			result.DockerCompose = string(fileContents)
		case "rancher-compose.tmpl":
			fileContents, err := getBytesFromURL(client, file.DownloadURL, token)
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
func (t *Template) SubBuildInfo(p *types.Plugin, tag string) error {
	val1, err1 := fixTemplate(p, "docker-compose.yml", t.DockerCompose, tag)
	if err1 != nil {
		return err1
	}
	t.DockerCompose = val1
	val2, err2 := fixTemplate(p, "rancher-compose.yml", t.RancherCompose, tag)
	if err2 != nil {
		return err2
	}
	t.RancherCompose = val2
	val3, err3 := fixTemplate(p, "config.yml", t.Config, tag)
	if err3 != nil {
		return err3
	}
	t.Config = val3
	return nil
}
