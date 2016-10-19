package github

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/h2non/gock"
)

func TestHookImage(t *testing.T) {
	defer gock.Off()
	g := goblin.Goblin(t)
	g.Describe("Github", func() {
		g.Describe("templateReplacement", func() {
			g.It("Check that files are being replaced with correct values", func() {
				var args tmplArguments
				args.Branch = "branch"
				args.Count = 1
				args.Tag = "tag"
				result, err := fixTemplate(&args, "docker-compose.yml", "{{ .Branch }} {{ .Count }} {{ .Tag }}")
				g.Assert(err).Equal(nil)
				g.Assert(result).Equal(fmt.Sprintf("%s %d %s", args.Branch, args.Count, args.Tag))
			})
		})
		g.Describe("Get data from url", func() {
			g.It("Should get 200", func() {
				defer gock.Off()
				gock.New("http://example.com").
					Get("/test").
					MatchHeader("Authorization", b64.URLEncoding.EncodeToString([]byte("token:x-oauth-basic"))).
					Reply(200).
					BodyString("hello!")
				client := &http.Client{}
				result, code, err := getBytesFromURL(client, "http://example.com/test", "token")
				g.Assert(err).Equal(nil)
				g.Assert(code).Equal(200)
				g.Assert(result).Equal([]byte("hello!"))
			})
		})
		g.Describe("get next template number", func() {
			g.It("should get the next number", func() {
				defer gock.Off()
				gock.New("http://example.com").
					Get("/test").
					MatchHeader("Authorization", b64.URLEncoding.EncodeToString([]byte("token:x-oauth-basic"))).
					Reply(200).
					JSON([]map[string]string{
					map[string]string{"name": "0"},
					map[string]string{"name": "1"},
					map[string]string{"name": "2"},
					map[string]string{"name": "SomeOtherFile"},
				})
				client := &http.Client{}
				number, err := getTemplateNum(client, "http://example.com/test", "token")
				g.Assert(err).Equal(nil)
				g.Assert(number).Equal(3)
			})

			g.It("should get 0 if there is no directory", func() {
				defer gock.Off()
				gock.New("http://example.com").
					Get("/test").
					MatchHeader("Authorization", b64.URLEncoding.EncodeToString([]byte("token:x-oauth-basic"))).
					Reply(401).JSON([]map[string]string{})
				client := &http.Client{}
				number, err := getTemplateNum(client, "http://example.com/test", "token")
				g.Assert(err).Equal(nil)
				g.Assert(number).Equal(0)
			})

		})
	})
}
