package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func main() {

}

// TODO: refactor funcs generations
// TODO: refactor structs


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
	commonPathTemp = "../output/"
)

type packages struct {
	name            string
	groupId         string
	artifactId      string
	destinationPath string
}

func init() {
	p := packages{
		name:            "mock-expedient",
		groupId:         "pe.mock.mock.expedient.app",
		artifactId:      "app",
		destinationPath: commonPathTemp,
	}

	// generateProject(p)
	generateBaseProject(p)
}

func generateProject(p packages) {
	generateOneProject(p)
}

func generateBaseProject(p packages) {
	projectPath := p.destinationPath + "mock-expedient"

	err := os.Mkdir(projectPath, 0755)
	if err != nil {
		panic(err)
	}

	pomPath := projectPath + "/pom.xml"
	generateRootPom(pomPath)

	multimoduleFlag := false
	hexagonalFlag := true

	if multimoduleFlag {
		generatePackages(p, projectPath)
		generatePackages2(p, projectPath)
	}

	if hexagonalFlag {
		generatePackages(p, projectPath)
		generatePackages2(p, projectPath)
		generatePackages3(p, projectPath)
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

func generatePackages(packages packages, projectpath string) {
	var secondLayer = []string{"main", "test"}
	var thirdLayer = []string{"java", "resources"}

	module := projectpath + "/mock-application"

	fmt.Print(module)

	err := os.Mkdir(module, 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		p := module + "/src"

		err := os.Mkdir(p, 0755)
		if err != nil {
			panic(err)
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}

		for j := 0; j < 2; j++ {
			p = absPath + "/" + secondLayer[j]

			err = os.Mkdir(p, 0755)
			if err != nil {
				panic(err)
			}

			newAbsPath, err := filepath.Abs(absPath + "/" + secondLayer[j])
			if err != nil {
				panic(err)
			}

			for k := 0; k < 2; k++ {
				main := newAbsPath + "/" + thirdLayer[k]

				err = os.Mkdir(main, 0755)
				if err != nil {
					panic(err)
				}

				if thirdLayer[k] == "java" {
					ns := main + "/pe.mock.expedient.app"
					err = os.Mkdir(ns, 0755)
					if err != nil {
						panic(err)
					}
				}
			}
		}
		pomPath := module + "/pom.xml"
		generatePom(pomPath)
		applicationPath := module + "/MockApplication.java"
		generateApplication(applicationPath)
	}
}

func generatePackages2(packages packages, projectpath string) {
	var secondLayer = []string{"main", "test"}
	var thirdLayer = []string{"java", "resources"}

	module := projectpath + "/mock-domain"

	fmt.Print(module)

	err := os.Mkdir(module, 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		p := module + "/src"

		err := os.Mkdir(p, 0755)
		if err != nil {
			panic(err)
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}

		for j := 0; j < 2; j++ {
			p = absPath + "/" + secondLayer[j]

			err = os.Mkdir(p, 0755)
			if err != nil {
				panic(err)
			}

			newAbsPath, err := filepath.Abs(absPath + "/" + secondLayer[j])
			if err != nil {
				panic(err)
			}

			for k := 0; k < 2; k++ {
				main := newAbsPath + "/" + thirdLayer[k]

				err = os.Mkdir(main, 0755)
				if err != nil {
					panic(err)
				}

				if thirdLayer[k] == "java" {
					ns := main + "/pe.mock.expedient.domain"
					err = os.Mkdir(ns, 0755)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		pomPath := module + "/pom.xml"
		generatePom(pomPath)
	}
}

func generatePackages3(packages packages, projectpath string) {
	var secondLayer = []string{"main", "test"}
	var thirdLayer = []string{"java", "resources"}

	module := projectpath + "/mock-infrastructure"

	fmt.Print(module)

	err := os.Mkdir(module, 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		p := module + "/src"

		err := os.Mkdir(p, 0755)
		if err != nil {
			panic(err)
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}

		for j := 0; j < 2; j++ {
			p = absPath + "/" + secondLayer[j]

			err = os.Mkdir(p, 0755)
			if err != nil {
				panic(err)
			}

			newAbsPath, err := filepath.Abs(absPath + "/" + secondLayer[j])
			if err != nil {
				panic(err)
			}

			for k := 0; k < 2; k++ {
				main := newAbsPath + "/" + thirdLayer[k]

				err = os.Mkdir(main, 0755)
				if err != nil {
					panic(err)
				}

				if thirdLayer[k] == "java" {
					ns := main + "/pe.mock.expedient.infra"
					err = os.Mkdir(ns, 0755)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		pomPath := module + "/pom.xml"
		generatePom(pomPath)
	}
}

func generateOneProject(packages packages) {
	var secondLayer = []string{"main", "test"}
	var thirdLayer = []string{"java", "resources"}

	module := packages.destinationPath + packages.name

	err := os.Mkdir(module, 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		p := module + "/src"

		err := os.Mkdir(p, 0755)
		if err != nil {
			panic(err)
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}

		for j := 0; j < 2; j++ {
			p = absPath + "/" + secondLayer[j]

			err = os.Mkdir(p, 0755)
			if err != nil {
				panic(err)
			}

			newAbsPath, err := filepath.Abs(absPath + "/" + secondLayer[j])
			if err != nil {
				panic(err)
			}

			for k := 0; k < 2; k++ {
				main := newAbsPath + "/" + thirdLayer[k]

				err = os.Mkdir(main, 0755)
				if err != nil {
					panic(err)
				}

				if thirdLayer[k] == "java" {
					ns := main + "/" + packages.groupId
					err = os.Mkdir(ns, 0755)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		pomPath := module + "/pom.xml"
		generatePom(pomPath)
		applicationPath := module + "/MockApplication.java"
		generateApplication(applicationPath)
	}
}
