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

// TODO: create new func for pom root
// TODO: create new func for main application class
// TODO: refactor structs
// TODO: refactor funcs generations

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

//go:embed templates/monorepo/pom.tmpl
var pomXmlTemplate embed.FS

//go:embed templates/multimodule/**
var x embed.FS

func generatePom(path string) {
	pomXmlTmpl := PomXmlTemplate{
		ArtifactIdParent: "ArtifactIdParentMOCK2",
		GroupIdParent:    "GroupIdParentMOCK2",
		VersionParent:    "0.181.0-1",
		ArtifactId:       "ArtifactIdMOCK2",
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

	tmpl, err := template.ParseFS(pomXmlTemplate, "templates/monorepo/pom.tmpl")
	if err != nil {
		fmt.Print("Failed to parse template")
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "pom.tmpl", pomXmlTmpl)
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
	// pomPath := p.destinationPath + p.name + "/pom.xml"
	// generatePom(pomPath)
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
	generatePom(pomPath)

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
	}
}
