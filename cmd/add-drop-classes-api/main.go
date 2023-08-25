package main

import (
	"github.com/alecthomas/kong"

	"github.com/TulgaCG/add-drop-classes-api/cmd/add-drop-classes-api/root"
)

func main() {
	ctx := kong.Parse(&root.Cli{},
		kong.Name("add-drop-classes-api"),
		kong.Description("A WebApp for Add/Drop Classes in a college"),
		kong.UsageOnError(),
	)

	ctx.FatalIfErrorf(ctx.Run())
}
