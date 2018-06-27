package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	yaml "gopkg.in/yaml.v2"
)

var InitCommand = cli.Command{
	Name: "init",
	Action: func(ctx *cli.Context) error {
		target := ctx.Args().First()
		ioutil.WriteFile(filepath.Join(target, "Appfile"), []byte("builder-image: golang"), 0755)

		return nil
	},
}

var BuildCommand = cli.Command{
	Name: "build",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "target, t",
			Usage: "build target",
		},
	},
	Action: func(ctx *cli.Context) error {
		buildDir := resolveBuildDir(ctx.Args().First())

		imageTag := ctx.String("target")
		if imageTag == "" {
			panic("image tag must not be empty")
		}

		appfile := loadAppfile(filepath.Join(buildDir, "Appfile"))

		build(buildDir, appfile.BuilderImage, appfile.Image, appfile.Bind, imageTag)

		return nil
	},
}

var RunCommand = cli.Command{
	Name: "run",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "command, c",
			Usage: "command to run",
		},
	},
	Action: func(ctx *cli.Context) error {
		buildDir := resolveBuildDir(ctx.Args().First())
		appfile := loadAppfile(filepath.Join(buildDir, "Appfile"))

		build(buildDir, appfile.BuilderImage, appfile.Image, appfile.Bind, "lol/wtf")

		command := ctx.String("command")
		if command == "" {
			command = appfile.Command
		}

		if command == "" {
			command = "/app"
		}

		// todo: remove lol/wtf and use --iid
		runCmd := exec.Command("docker", "run", "--rm", "-i", "lol/wtf", command)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			panic(err)
		}

		return nil
	},
}

var TestCommand = cli.Command{
	Name: "test",
	Action: func(ctx *cli.Context) error {
		buildDir := resolveBuildDir(ctx.Args().First())
		appfile := loadAppfile(filepath.Join(buildDir, "Appfile"))

		build(buildDir, appfile.BuilderImage, appfile.Image, appfile.Bind, "lol/wtf")

		runCmd := exec.Command("docker", "run", "--rm", "-i", "lol/wtf", appfile.Test)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			panic(err)
		}

		return nil
	},
}

type Appfile struct {
	BuilderImage string `yaml:"builder-image"`
	Image        string
	Bind         string
	Command      string
	Test         string
}

func loadAppfile(appfilePath string) *Appfile {
	file, err := os.Open(appfilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var appfile Appfile
	if err := yaml.NewDecoder(file).Decode(&appfile); err != nil {
		panic(err)
	}

	return &appfile
}

func build(buildDir, builderImage, image, bind, tag string) {
	dockerfile := fmt.Sprintf("FROM %s\nADD . %s", image, bind)
	if builderImage == "golang" {
		dockerfile = `
		FROM golang:1.8 as build

		WORKDIR /go/src/app
		COPY . .

		RUN go-wrapper download   # "go get -d -v ./..."
		RUN go-wrapper install

		FROM busybox
		COPY --from=build /go/bin/app /
		CMD ["/app"]
		`
	}

	buildCmd := exec.Command("docker", "build", buildDir, "-t", tag, "-f", "-")
	buildCmd.Stdin = bytes.NewBuffer([]byte(dockerfile))
	if output, err := buildCmd.CombinedOutput(); err != nil {
		panic(string(output))
	}
}

func resolveBuildDir(dir string) string {
	u, err := url.Parse(dir)
	if err != nil {
		panic(err)
	}

	if u.Scheme == "https" {
		tmp, err := ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}

		if out, err := exec.Command("git", "clone", dir, tmp).CombinedOutput(); err != nil {
			panic(string(out))
		}

		return tmp
	}

	return dir
}

func main() {
	app := cli.NewApp()
	app.Name = "runa"
	app.Commands = []cli.Command{
		BuildCommand,
		RunCommand,
		InitCommand,
		TestCommand,
	}

	app.Run(os.Args)
}
