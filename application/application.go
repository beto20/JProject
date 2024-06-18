package application

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/beto20/jproject/util"
)

type Input struct {
	Group       string
	Artifact    string
	Name        string
	Description string
	PackageName string
	JavaVersion string
	DestinyPath string
	ProjectType string
	Prefix      string
	Module      []Module
}

type PomRootXmlTemplate struct {
	ArtifactIdParent string
	GroupIdParent    string
	VersionParent    string
	ArtifactId       string
	GroupId          string
	Version          string
	Module           []Module
	JavaVersion      int64
	PomDepTmpl       []PomDependencyTemplate
	Scm              Scm
	Repository       Repository
}
type Module struct {
	Name       string
	ArtifactId string
}
type Scm struct {
	HasScm  bool
	Project string
}
type Repository struct {
	HasRepository bool
	Id            string
	Name          string
	Url           string
}

type PomXmlTemplate struct { //PomTemplate
	ArtifactIdParent string
	GroupIdParent    string
	VersionParent    string
	ArtifactId       string
	GroupId          string
	PomDepTmpl       []PomDependencyTemplate
}
type PomDependencyTemplate struct { // DependencyTemplate
	GroupIdDependency    string
	ArtifactIdDependency string
}

type ApplicationTemplate struct { // MainApplicationTemplate
	Namespace    string
	BasePackages string
	Name         string
}

type config struct { // baseConfig
	modu               string
	namespace          string
	artifact           string
	requireApplication bool
}

//go:embed templates/monorepo/pom.xml
var pomXmlTemplate embed.FS

//go:embed templates/multimodule/pom_root.xml
var pomRootXmlTemplate embed.FS

//go:embed templates/main_application.tmpl
var applicationTemplate embed.FS

func toApplicationTemplate(i Input) ApplicationTemplate {
	// mainClassName := strings.ReplaceAll(i.Artifact, "-", "")

	parts := strings.Split(i.PackageName, ".")

	base := "mock"

	if len(parts) >= 2 {
		base = parts[0] + "." + parts[1]
	} else if len(parts) == 1 {
		base = parts[0]
	} else {
		fmt.Println("The string does not contain at least two words.")
	}

	appTmpl := ApplicationTemplate{
		Namespace:    i.PackageName,
		BasePackages: base,
		Name:         i.Artifact,
	}

	return appTmpl
}

func ToPomXmlTemplate(i Input, index int64) PomXmlTemplate {

	artifact := i.Artifact

	if i.ProjectType == "hexagonal" {
		artifact = i.Module[index].ArtifactId
	} else if i.ProjectType == "multimodule" {
		artifact = i.Module[index].ArtifactId
	}

	pomTmpl := PomXmlTemplate{
		// TODO: replace this values
		ArtifactIdParent: util.ARTIFACT_ID_PARENT,
		GroupIdParent:    util.GROUP_ID_PARENT,
		VersionParent:    util.VERSION_PARENT,
		ArtifactId:       artifact,
		GroupId:          i.Group,
	}

	return pomTmpl
}

func ToPomRootXmlTemplate(i Input) PomRootXmlTemplate {
	jv, err := strconv.ParseInt(i.JavaVersion, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	pomRootXmlTmpl := PomRootXmlTemplate{
		ArtifactIdParent: util.ARTIFACT_ID_PARENT,
		GroupIdParent:    util.GROUP_ID_PARENT,
		VersionParent:    util.VERSION_PARENT,

		ArtifactId:  i.Artifact,
		GroupId:     i.Group,
		Version:     "1.0.0",
		Module:      i.Module,
		JavaVersion: jv,
		Scm: Scm{
			HasScm:  true,
			Project: util.SCM_PROJECT + i.Artifact,
		},
		Repository: Repository{
			HasRepository: true,
			Id:            util.REPOSITORY_ID,
			Name:          util.REPOSITORY_NAME,
			Url:           util.REPOSITORY_URL,
		},
	}

	return pomRootXmlTmpl
}

func GenerateProject(input Input) {
	config := setInitProjectConfig(input)
	input = toModules(config, input)

	fmt.Println("input", input)
	var projectPath = ""

	if input.ProjectType == "module" {
		projectPath = input.DestinyPath
	} else {
		projectPath = input.DestinyPath + input.Artifact
		err := os.Mkdir(projectPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	generate(projectPath, config, input)

	if input.ProjectType != "monorepo" && input.ProjectType != "module" {
		pomPath := projectPath + "/pom.xml"
		generateRootPom(pomPath, input)
	}
}

func toModules(configs []config, input Input) Input {
	var mods []Module

	for _, c := range configs {
		mod := Module{
			Name:       strings.ReplaceAll(c.modu, "/", ""),
			ArtifactId: c.artifact,
		}

		mods = append(mods, mod)
	}

	input.Module = mods

	return input
}

func setInitProjectConfig(i Input) []config {
	var bc []config

	ns := "/" + i.PackageName
	px := "/" + i.Prefix
	n := "/" + i.Name

	if i.ProjectType == "hexagonal" {
		bc = []config{
			{
				modu:               px + "-application",
				namespace:          ns + ".app",
				artifact:           i.Artifact + "-app",
				requireApplication: false,
			},
			{
				modu:               px + "-domain",
				namespace:          ns + ".domain",
				artifact:           i.Artifact + "-domain",
				requireApplication: false,
			},
			{
				modu:               px + "-infrastructure",
				namespace:          ns + ".infra",
				artifact:           i.Artifact + "-infra",
				requireApplication: true,
			},
		}
	}

	if i.ProjectType == "multimodule" {
		bc = []config{
			{
				modu:               px + "-app",
				namespace:          ns + ".app",
				artifact:           i.Artifact + "-app",
				requireApplication: true,
			},
			{
				modu:               px + "-core",
				namespace:          ns + ".core",
				artifact:           i.Artifact + "-core",
				requireApplication: false,
			},
		}
	}

	if i.ProjectType == "monorepo" {
		bc = []config{
			{
				modu:               "",
				namespace:          ns,
				artifact:           "",
				requireApplication: true,
			},
		}
	}

	if i.ProjectType == "module" {
		bc = []config{
			{
				modu:               n,
				namespace:          ns,
				artifact:           "",
				requireApplication: false,
			},
		}
	}

	return bc
}

func generate(projectPath string, project []config, input Input) {
	var index int64 = 0
	for _, p := range project {
		generatePackages(projectPath, p)
		pom(projectPath, p.modu, input, index)
		app(projectPath, p, input)
		index++
	}
}

func generatePom(path string, input Input, index int64) {
	// TODO: remove this when start refactoring
	// pomXmlTmpl := PomXmlTemplate{
	// 	ArtifactIdParent: "ArtifactIdParentMOCK2",
	// 	GroupIdParent:    "GroupIdParentMOCK2",
	// 	VersionParent:    "0.181.0-1",
	// 	ArtifactId:       "ArtifactIdMOCK2",
	// 	GroupId:          "GroupIdMOCK2",
	// PomDepTmpl: []PomDependencyTemplate{
	// 	{
	// 		GroupIdDependency:    "org.springframework.boot",
	// 		ArtifactIdDependency: "spring-boot-starter-data-jpa",
	// 	},
	// 	{
	// 		GroupIdDependency:    "org.2",
	// 		ArtifactIdDependency: "spring-boot-starter-3",
	// 	},
	// },
	// }

	pomXmlTmpl := ToPomXmlTemplate(input, index)

	tmpl, err := template.ParseFS(pomXmlTemplate, "templates/monorepo/pom.xml")
	if err != nil {
		fmt.Print("Failed to parse template")
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "pom.xml", pomXmlTmpl)
	if err != nil {
		fmt.Print("Failed to execute template")
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, pomXmlTmpl)
	if err != nil {
		fmt.Print("Failed to execute template, generate output file")
	}
}

func generateRootPom(path string, input Input) {
	// TODO: remove this when start refactoring
	// pomRootXmlTmpl := PomRootXmlTemplate{
	// 	ArtifactIdParent: "ArtifactIdParentMOCK2",
	// 	GroupIdParent:    "GroupIdParentMOCK2",
	// 	VersionParent:    "0.181.0-1",
	// 	ArtifactId:       "ArtifactIdMOCK2",
	// 	GroupId:          "GroupIdMOCK2",
	// 	Version:          "1.0.0",
	// 	Module: []Module{
	// 		{Name: "mock-app"},
	// 		{Name: "mock-core"},
	// 	},
	// 	JavaVersion: 11,
	// 	PomDepTmpl: []PomDependencyTemplate{
	// 		{
	// 			GroupIdDependency:    "org.springframework.boot",
	// 			ArtifactIdDependency: "spring-boot-starter-data-jpa",
	// 		},
	// 		{
	// 			GroupIdDependency:    "org.2",
	// 			ArtifactIdDependency: "spring-boot-starter-3",
	// 		},
	// 	},
	// 	Scm: Scm{
	// 		HasScm:  true,
	// 		Project: "bitbucket.mock-expedient",
	// 	},
	// 	Repository: Repository{
	// 		HasRepository: true,
	// 		Id:            "central",
	// 		Name:          "central mock",
	// 		Url:           "https://central-mock.com",
	// 	},
	// }

	pomRootXmlTmpl := ToPomRootXmlTemplate(input)

	tmpl, err := template.ParseFS(pomRootXmlTemplate, "templates/multimodule/pom_root.xml")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, pomRootXmlTmpl)
	if err != nil {
		panic(err)
	}
}

func generateApplication(path string, input Input) {
	// TODO: remove this when start refactoring
	// applicationTmpl := ApplicationTemplate{
	// 	Namespace:    "pe.com.mock",
	// 	BasePackages: "pe.com",
	// 	Name:         "MockApplication",
	// }

	applicationTmpl := toApplicationTemplate(input)

	tmpl, err := template.ParseFS(applicationTemplate, "templates/main_application.tmpl")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, applicationTmpl)
	if err != nil {
		panic(err)
	}
}

func generatePackages(projectpath string, c config) {
	module := projectpath + c.modu

	if c.modu != "" {
		err := os.Mkdir(module, 0755)
		if err != nil {
			panic(err)
		}
	}

	p := module + "/src"

	err := os.Mkdir(p, 0755)
	if err != nil {
		panic(err)
	}

	buildSrc(p, c)
}

func pom(projectpath string, c string, input Input, index int64) {
	module := projectpath + c
	pomPath := module + "/pom.xml"
	generatePom(pomPath, input, index)
}

func app(projectpath string, c config, input Input) {
	module := projectpath + c.modu

	if c.requireApplication {
		input.Artifact = strings.ReplaceAll(input.Artifact, "-", "")
		applicationPath := module + "/src/main/java/" + c.namespace + "/" + input.Artifact
		generateApplication(applicationPath, input)
	}
}

func buildSrc(p string, c config) {
	var firstLayer = []string{"main", "test"}
	var secondLayer = []string{"java", "resources"}

	absPath, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i++ {
		p = absPath + "/" + firstLayer[i]

		err = os.Mkdir(p, 0755)
		if err != nil {
			panic(err)
		}

		newAbsPath, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}

		for j := 0; j < 2; j++ {
			np := newAbsPath + "/" + secondLayer[j]

			err = os.Mkdir(np, 0755)
			if err != nil {
				panic(err)
			}

			if secondLayer[j] == "java" {
				ns := np + c.namespace
				err = os.Mkdir(ns, 0755)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func help() {

	helpCommand := `
	*************************************************************************************
	*         					    	JProject										*
	*************************************************************************************
	*         					    	Subcommand										*
	* help				-h		show all commands descriptions   						*
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
	
	supported by av
	`

	fmt.Println(helpCommand)
}
