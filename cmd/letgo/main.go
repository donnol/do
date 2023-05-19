package main

import (
	"context"
	"os"

	"github.com/donnol/do"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Commands = cmds

	do.Must(app.RunContext(context.Background(), os.Args))
}
