package docker

import (
	"testing"

	"github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Docker Command Creation", func() {
		g.It("Command should have correct arguments", func() {
			g.Assert(CreateCmd([]string{"abc", "efg"}, true).Args).Equal([]string{DockerCmd, "abc", "efg"})
		})
	})
}
