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
	Run(sc string) error
}

var availableCommand = [6]string{"h", "v", "mm", "mr", "mod", "hex"}

func NewCommand(cm string) *Command {
	cmd := &Command{
		fs: flag.NewFlagSet(cm, flag.ContinueOnError),
	}

	return genericFlags(cmd)
}

func genericFlags(c *Command) *Command {
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

func (c *Command) Run(sc string) error {

	switch sc {
	case "h":
		help()
	case "v":
		version()
	default:
		create(c, sc)
	}

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

	res := constainsCommand(availableCommand, subcommand)

	if !res {
		return fmt.Errorf("Unknown subcommand: %s", subcommand)
	}

	cmds := []Runner{
		NewCommand(subcommand),
	}

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run(os.Args[1])
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func constainsCommand(arr [6]string, value string) bool {
	for _, c := range arr {
		if c == value {
			return true
		}
	}
	return false
}

func create(c *Command, sc string) {
	var projectType string = ""

	switch sc {
	case "hex":
		projectType = "hexagonal"
	case "mm":
		projectType = "multimodule"
	case "mr":
		projectType = "monorepo"
	case "mod":
		projectType = "module"
	}

	i := input{
		group:       c.group,
		artifact:    c.artifact,
		name:        c.name,
		description: c.description,
		packageName: c.packageName,
		javaVersion: c.javaVersion,
		destinyPath: c.destinyPath,
		projectType: projectType,
		prefix:      c.prefix,
	}

	application.GenerateProject(toApplicationInput(i))
}

func help() {

	helpCommand := `
	*************************************************************************************
	*         					    	JProject										*
	*************************************************************************************
	*         					    	Subcommand										*
	* help				-h		show all commands descriptions   						*
	* version			-v		show the current version		  						*
	* monorepo			-mr		generate a monorepo project   							*
	* module			-m		generate an empty module    							*
	* multimodule		-mm		generate a multimodule project (app and core mods)   	*
	* hexagonal			-hex	generate a hexagonal project (infra, app and domain)   	*
	*************************************************************************************
	*         					    	Flags											*
	* Important advice: for each generate subcommand you must specify any flag			*
	* groupId			-g		project groupId					   		requiered		*
	* artifactId		-a		project artifactId   					requiered		*
	* name				-n		project name   							optional		*
	* description		-d		project description   					optional		*
	* package			-pk		project package   						requiered		*
	* project prefix	-p		project prefix just for mm and hex 		requiered		*
	* java version		-jv		project java version (8 or 11)   		requiered		*
	* directory path	-dp		ubication of project generated   		requiered		*
	*************************************************************************************
	`

	fmt.Println(helpCommand)
}

func version() {
	v := "JProject v1.0.0 2024-06-20"
	fmt.Println(v)
}
