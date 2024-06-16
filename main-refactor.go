// package main

// import (
// 	"embed"
// 	_ "embed"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"text/template"
// )

// func mainref() {

// }

// type input struct {
// 	group       string
// 	artifact    string
// 	name        string
// 	description string
// 	packageName string
// 	javaVersion string
// 	destinyPath string
// 	projectType string
// 	prefix      string
// }

// type MainApplicationTemplate struct {
// 	Namespace       string
// 	BasePackages    string
// 	Name            string
// 	pathDestination string
// }

// type PomTemplate struct {
// 	ArtifactIdParent string
// 	GroupIdParent    string
// 	VersionParent    string
// 	ArtifactId       string
// 	GroupId          string
// 	PomDepTmpl       []DependencyTemplate
// 	pathDestination  string
// }

// type DependencyTemplate struct {
// 	GroupIdDependency    string
// 	ArtifactIdDependency string
// }

// type baseConfig struct {
// 	modu               string
// 	namespace          string
// 	requireApplication bool
// }

// type pckg struct {
// 	name            string
// 	groupId         string
// 	artifactId      string
// 	destinationPath string
// }

// //go:embed templates/main_application.tmpl
// var mainAppTemplate embed.FS

// //go:embed templates/monorepo/pom.xml
// var pomTemplate embed.FS

// func mainFunc() {

// 	// mappers
// 	mainTemplate := transformInputToTemplates(i)
// 	pomTemplate := transformInputToPomTemplate(i)

// 	// setup init configuration based on project type
// 	config := setInitProjectConfig(i)

// 	// project structure generation
// 	create(config)

// 	// templates builds
// 	buildMainApplication(mainTemplate)
// 	buildPom(pomTemplate)
// }

// func create(config []baseConfig) {

// 	x := toPckg()
// 	buildModule(x)
// }

// func toPckg() pckg {
// 	p := pckg{
// 		name:            "mock-expedient",
// 		groupId:         "pe.mock.mock.expedient.app",
// 		artifactId:      "app",
// 		destinationPath: commonPathTemp,
// 	}

// 	return p
// }

// func buildModule(p pckg) {
// 	var secondLayer = []string{"main", "test"}
// 	var thirdLayer = []string{"java", "resources"}

// 	module := p.destinationPath + "mock-expedient" + p.mod

// 	err := os.Mkdir(module, 0755)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < 1; i++ {
// 		p := module + "/src"

// 		err := os.Mkdir(p, 0755)
// 		if err != nil {
// 			panic(err)
// 		}
// 		absPath, err := filepath.Abs(p)
// 		if err != nil {
// 			panic(err)
// 		}

// 		for j := 0; j < 2; j++ {
// 			p = absPath + "/" + secondLayer[j]

// 			err = os.Mkdir(p, 0755)
// 			if err != nil {
// 				panic(err)
// 			}

// 			newAbsPath, err := filepath.Abs(absPath + "/" + secondLayer[j])
// 			if err != nil {
// 				panic(err)
// 			}

// 			for k := 0; k < 2; k++ {
// 				main := newAbsPath + "/" + thirdLayer[k]

// 				err = os.Mkdir(main, 0755)
// 				if err != nil {
// 					panic(err)
// 				}

// 				if thirdLayer[k] == "java" {
// 					ns := main + namespace
// 					err = os.Mkdir(ns, 0755)
// 					if err != nil {
// 						panic(err)
// 					}
// 				}
// 			}
// 		}

// 		pomPath := module + "/pom.xml"
// 		generatePom(pomPath)

// 		if requireApplication {
// 			applicationPath := module + "/MockApplication.java"
// 			generateApplication(applicationPath)
// 		}
// 	}
// }

// func setInitProjectConfig(i input) []baseConfig {
// 	var bc []baseConfig

// 	ns := "/" + i.packageName
// 	px := "/" + i.prefix
// 	n := "/" + i.name

// 	if i.projectType == "monorepo" {
// 		bc = []baseConfig{
// 			{
// 				modu:               n,
// 				namespace:          ns,
// 				requireApplication: true,
// 			},
// 		}
// 	}

// 	if i.projectType == "module" {
// 		bc = []baseConfig{
// 			{
// 				modu:               n,
// 				namespace:          ns,
// 				requireApplication: false,
// 			},
// 		}
// 	}

// 	if i.projectType == "multimodule" {
// 		bc = []baseConfig{
// 			{
// 				modu:               px + "-app",
// 				namespace:          ns + ".app",
// 				requireApplication: true,
// 			},
// 			{
// 				modu:               px + "-core",
// 				namespace:          ns + ".core",
// 				requireApplication: false,
// 			},
// 		}
// 	}

// 	if i.projectType == "hexagonal" {
// 		bc = []baseConfig{
// 			{
// 				modu:               px + "-application",
// 				namespace:          ns + ".app",
// 				requireApplication: false,
// 			},
// 			{
// 				modu:               px + "-domain",
// 				namespace:          ns + ".domain",
// 				requireApplication: false,
// 			},
// 			{
// 				modu:               px + "-infrastructure",
// 				namespace:          ns + ".infra",
// 				requireApplication: true,
// 			},
// 		}
// 	}

// 	return bc
// }

// func transformInputToTemplates(i input) MainApplicationTemplate {
// 	mainClassName := strings.ReplaceAll(i.artifact, "-", "")

// 	appTmpl := MainApplicationTemplate{
// 		Namespace:       i.packageName,
// 		BasePackages:    "pe.xxxx", // TODO
// 		Name:            mainClassName,
// 		pathDestination: "", // TODO to validate
// 	}

// 	return appTmpl
// }

// func transformInputToPomTemplate(i input) PomTemplate {

// 	pomTmpl := PomTemplate{
// 		// TODO: replace this values
// 		ArtifactIdParent: "ArtifactIdParentMOCK2",
// 		GroupIdParent:    "GroupIdParentMOCK2",
// 		VersionParent:    "0.181.0-1",

// 		ArtifactId: i.artifact,
// 		GroupId:    i.group,
// 		PomDepTmpl: []DependencyTemplate{ // TODO improve dependencies builder
// 			{
// 				GroupIdDependency:    "org.springframework.boot",
// 				ArtifactIdDependency: "spring-boot-starter-data-jpa",
// 			},
// 			{
// 				GroupIdDependency:    "org.2",
// 				ArtifactIdDependency: "spring-boot-starter-3",
// 			},
// 		},
// 		pathDestination: "", // TODO to validate
// 	}

// 	return pomTmpl
// }

// func buildMainApplication(mainTmpl MainApplicationTemplate) {
// 	tmpl, err := template.ParseFS(mainAppTemplate, "templates/main_application.tmpl")
// 	if err != nil {
// 		panic(err)
// 	}

// 	file, err := os.Create(mainTmpl.pathDestination)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = tmpl.Execute(file, mainTmpl)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func buildPom(pomTmpl PomTemplate) {
// 	tmpl, err := template.ParseFS(pomTemplate, "templates/monorepo/pom.xml")
// 	if err != nil {
// 		fmt.Print("Failed to parse template")
// 	}

// 	err = tmpl.ExecuteTemplate(os.Stdout, "pom.xml", pomTmpl)
// 	if err != nil {
// 		fmt.Print("Failed to execute template")
// 	}

// 	file, err := os.Create(pomTmpl.pathDestination)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = tmpl.Execute(file, pomTmpl)
// 	if err != nil {
// 		fmt.Print("Failed to execute template, generate output file")
// 	}
// }
