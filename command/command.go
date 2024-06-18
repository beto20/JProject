package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/beto20/jproject/application"
)

type Command struct {
	fs          *flag.FlagSet
	group       string
	artifact    string
	name        string
	description string
	packageName string
	javaVersion string
	destinyPath string
	prefix      string
}

type input struct {
	group       string
	artifact    string
	name        string
	description string
	packageName string
	javaVersion string
	destinyPath string
	projectType string
	prefix      string
}

type Runner interface {
	Init([]string) error
	Name() string
	Run() error
}

func NewCommand() *Command {

	c := &Command{
		fs: flag.NewFlagSet("hex", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.group, "g", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.artifact, "a", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.name, "n", "DEFAULT_TEST", "description")
	c.fs.StringVar(&c.description, "d", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.packageName, "pk", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.javaVersion, "jv", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.destinyPath, "dp", "DEFAULT_TEST", "")
	c.fs.StringVar(&c.prefix, "p", "DEFAULT_TEST", "")

	return c
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Run() error {

	// i := input{
	// 	group:       c.group,
	// 	artifact:    c.artifact,
	// 	name:        c.name,
	// 	description: c.description,
	// 	packageName: c.packageName,
	// 	javaVersion: c.javaVersion,
	// 	destinyPath: c.destinyPath,
	// 	projectType: "multimodule",
	// 	prefix:      c.prefix,
	// }

	// application.GenerateProject(toApplicationInput(i))

	return nil
}

func toApplicationInput(i input) application.Input {
	appInput := application.Input{
		Group:       i.group,
		Artifact:    i.artifact,
		Name:        i.name,
		Description: i.description,
		PackageName: i.packageName,
		JavaVersion: i.javaVersion,
		DestinyPath: i.destinyPath,
		ProjectType: i.projectType,
		Prefix:      i.prefix,
	}

	return appInput
}

func Root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	subcommand := os.Args[1]

	cmds := []Runner{
		NewCommand(),
	}

	fmt.Println("subcommand:", subcommand)

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			fmt.Println("cmd:", os.Args[2:])

			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
