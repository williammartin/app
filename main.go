package main

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

var BuildCommand = cli.Command{
	Name: "build",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "build target",
		},
	},
	Action: func(ctx *cli.Context) error {
		buildDir := ctx.Args().First()

		cmd := exec.Command("docker", "build", buildDir, "-t", "my/app", "-f", "-")
		cmd.Stdin = bytes.NewBuffer([]byte("FROM scratch\nADD . tmp/app"))
		if output, err := cmd.CombinedOutput(); err != nil {
			panic(string(output))
		}

		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "app"
	app.Commands = []cli.Command{
		BuildCommand,
	}

	app.Run(os.Args)
}
