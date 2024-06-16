package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Command struct {
	fs *flag.FlagSet

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
	Run() error
	Name() string
}

func NewCommand() *Command {
	gc := &Command{
		fs: flag.NewFlagSet("mm", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.group, "g", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.artifact, "a", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.name, "n", "DEFAULT_TEST", "description")
	gc.fs.StringVar(&gc.description, "d", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.packageName, "pk", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.javaVersion, "jv", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.destinyPath, "dp", "DEFAULT_TEST", "")
	gc.fs.StringVar(&gc.prefix, "p", "DEFAULT_TEST", "")

	return gc
}

func (g *Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Command) Run() error {

	i := input{
		group:       g.group,
		artifact:    g.artifact,
		name:        g.name,
		description: g.description,
		packageName: g.packageName,
		javaVersion: g.javaVersion,
		destinyPath: g.destinyPath,
		projectType: "monorepo",
		prefix:      g.prefix,
	}

	fmt.Println("-g", i.group)
	fmt.Println("-a", i.artifact)
	fmt.Println("-n", i.name)
	fmt.Println("-d", i.description)
	fmt.Println("-pk", i.packageName)
	fmt.Println("-jv", i.javaVersion)
	fmt.Println("-dp", i.destinyPath)
	fmt.Println("-p", i.prefix)

	// TODO: integrate with the project generator core

	return nil
}

func (g *Command) Name() string {
	return g.fs.Name()
}

func Root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		NewCommand(),
	}
	fmt.Println("cmds:", cmds)

	subcommand := os.Args[1]

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
