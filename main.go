package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/beto20/jproject/command"
)

// TODO: define input command line arguments

type inputProject struct {
	group       string
	artifact    string
	name        string
	description string
	packageName string
	javaVersion string
	destinyPath string
	projectType string
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
	Name string
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

type PomXmlTemplate struct {
	ArtifactIdParent string
	GroupIdParent    string
	VersionParent    string
	ArtifactId       string
	GroupId          string
	PomDepTmpl       []PomDependencyTemplate
}
type PomDependencyTemplate struct {
	GroupIdDependency    string
	ArtifactIdDependency string
}

type ApplicationTemplate struct {
	Namespace    string
	BasePackages string
	Name         string
}

//go:embed templates/monorepo/pom.xml
var pomXmlTemplate embed.FS

//go:embed templates/multimodule/pom_root.xml
var pomRootXmlTemplate embed.FS

//go:embed templates/main_application.tmpl
var applicationTemplate embed.FS

//go:embed templates/multimodule/**
var x embed.FS

const (
	commonPathTemp = ".../output/"
)

type packages struct {
	name            string
	groupId         string
	artifactId      string
	destinationPath string
}

type config struct {
	modu               string
	namespace          string
	requireApplication bool
}

func main() {

	if err := command.Root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init2() {
	p := packages{
		name:            "mock-expedient1",
		groupId:         "pe.mock.mock.expedient.app",
		artifactId:      "app",
		destinationPath: commonPathTemp,
	}

	input := "multimodule"
	generateProject(p, input)
}

func generateProject(packages packages, input string) {
	config := setProjectConfiguration(input)
	var projectPath = ""

	if input == "module" {
		projectPath = packages.destinationPath
	} else {
		projectPath = packages.destinationPath + packages.name
		err := os.Mkdir(projectPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	generate(projectPath, config)

	if input != "monorepo" && input != "module" {
		pomPath := projectPath + "/pom.xml"
		generateRootPom(pomPath)
	}
}

func setProjectConfiguration(input string) []config {
	var c []config

	if input == "hexagonal" {
		c = []config{
			{
				modu:               "/mock-application",
				namespace:          "/pe.mock.expedient.app",
				requireApplication: false,
			},
			{
				modu:               "/mock-domain",
				namespace:          "/pe.mock.expedient.domain",
				requireApplication: false,
			},
			{
				modu:               "/mock-infrastructure",
				namespace:          "/pe.mock.expedient.infra",
				requireApplication: true,
			},
		}
	}

	if input == "multimodule" {
		c = []config{
			{
				modu:               "/mock-app",
				namespace:          "/pe.mock.expedient.app",
				requireApplication: true,
			},
			{
				modu:               "/mock-core",
				namespace:          "/pe.mock.expedient.core",
				requireApplication: false,
			},
		}
	}

	if input == "monorepo" {
		c = []config{
			{
				modu:               "",
				namespace:          "/pe.mock.expedient",
				requireApplication: true,
			},
		}
	}

	if input == "module" {
		c = []config{
			{
				modu:               "new-mod-3",
				namespace:          "/pe.mock.expedient",
				requireApplication: false,
			},
		}
	}

	return c
}

func generate(projectPath string, project []config) {
	for _, p := range project {
		generatePackages(projectPath, p)
		pom(projectPath, p.modu)
		app(projectPath, p)
	}
}

func generatePom(path string) {
	pomXmlTmpl := PomXmlTemplate{
		ArtifactIdParent: "ArtifactIdParentMOCK2",
		GroupIdParent:    "GroupIdParentMOCK2",
		VersionParent:    "0.181.0-1",
		ArtifactId:       "ArtifactIdMOCK2",
		GroupId:          "GroupIdMOCK2",
		PomDepTmpl: []PomDependencyTemplate{
			{
				GroupIdDependency:    "org.springframework.boot",
				ArtifactIdDependency: "spring-boot-starter-data-jpa",
			},
			{
				GroupIdDependency:    "org.2",
				ArtifactIdDependency: "spring-boot-starter-3",
			},
		},
	}

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

func generateRootPom(path string) {

	pomRootXmlTmpl := PomRootXmlTemplate{
		ArtifactIdParent: "ArtifactIdParentMOCK2",
		GroupIdParent:    "GroupIdParentMOCK2",
		VersionParent:    "0.181.0-1",
		ArtifactId:       "ArtifactIdMOCK2",
		GroupId:          "GroupIdMOCK2",
		Version:          "1.0.0",
		Module: []Module{
			{Name: "mock-app"},
			{Name: "mock-core"},
		},
		JavaVersion: 11,
		PomDepTmpl: []PomDependencyTemplate{
			{
				GroupIdDependency:    "org.springframework.boot",
				ArtifactIdDependency: "spring-boot-starter-data-jpa",
			},
			{
				GroupIdDependency:    "org.2",
				ArtifactIdDependency: "spring-boot-starter-3",
			},
		},
		Scm: Scm{
			HasScm:  true,
			Project: "bitbucket.mock-expedient",
		},
		Repository: Repository{
			HasRepository: true,
			Id:            "central",
			Name:          "central mock",
			Url:           "https://central-mock.com",
		},
	}

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

func generateApplication(path string) {
	applicationTmpl := ApplicationTemplate{
		Namespace:    "pe.com.mock",
		BasePackages: "pe.com",
		Name:         "MockApplication",
	}

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

func pom(projectpath string, c string) {
	module := projectpath + c
	pomPath := module + "/pom.xml"
	generatePom(pomPath)
}

func app(projectpath string, c config) {
	module := projectpath + c.modu
	if c.requireApplication {
		applicationPath := module + "/src/main/java/" + c.namespace + "/MockApplication.java"
		generateApplication(applicationPath)
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
