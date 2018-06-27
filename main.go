package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	yaml "gopkg.in/yaml.v2"
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

		imageTag := ctx.String("target")
		if imageTag == "" {
			panic("image tag must not be empty")
		}

		file, err := os.Open(filepath.Join(buildDir, "Appfile"))
		if err != nil {
			panic(err)
		}

		defer file.Close()

		manifest := struct {
			Image string
			Bind  string
		}{}
		if err := yaml.NewDecoder(file).Decode(&manifest); err != nil {
			panic(err)
		}

		dockerfile := fmt.Sprintf("FROM %s\nADD . %s", manifest.Image, manifest.Bind)
		cmd := exec.Command("docker", "build", buildDir, "-t", imageTag, "-f", "-")
		cmd.Stdin = bytes.NewBuffer([]byte(dockerfile))
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
