package root

import (
	"github.com/TulgaCG/add-drop-classes-api/cmd/add-drop-classes-api/server"
)

type Cli struct {
	Server server.Cmd `cmd:"" help:"Run the server"`
}
