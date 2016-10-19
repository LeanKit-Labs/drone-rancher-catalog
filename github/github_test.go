package github

import (
	"fmt"
	"testing"

	"github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Github", func() {
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
}
