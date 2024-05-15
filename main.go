package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

func main() {
}

type xmlTemplate struct {
	parent struct {
		artifactIdParent string
		groupIdParent    string
		versionParent    string
	}
	artifactId string
	dependency struct {
		groupIdDependency    string
		artifactIdDependency string
	}
}

type XmlTemplate2 struct {
	ArtifactIdParent     string
	GroupIdParent        string
	VersionParent        string
	ArtifactId           string
	GroupIdDependency    string
	ArtifactIdDependency string
}

//go:embed templates/pom_template.xml
var pomXmlTemplate embed.FS

func init() {
	fmt.Println("init")

	// xmlTemplate := xmlTemplate{
	// 	parent: struct {
	// 		artifactIdParent string
	// 		groupIdParent    string
	// 		versionParent    string
	// 	}{
	// 		artifactIdParent: "ArtifactId",
	// 		groupIdParent:    "groupIdParent",
	// 		versionParent:    "0.181.0-1-SNAPSHOT",
	// 	},
	// 	artifactId: "artifactId",
	// 	dependency: struct {
	// 		groupIdDependency    string
	// 		artifactIdDependency string
	// 	}{
	// 		groupIdDependency:    "org.springframework.boot",
	// 		artifactIdDependency: "spring-boot-starter-data-jpa",
	// 	},
	// }

	xmlTemplate2 := XmlTemplate2{
		ArtifactIdParent:     "ArtifactIdParent",
		GroupIdParent:        "GroupIdParent",
		VersionParent:        "0.181.0-1",
		ArtifactId:           "ArtifactId",
		GroupIdDependency:    "org.springframework.boot",
		ArtifactIdDependency: "spring-boot-starter-data-jpa",
	}

	tmpl, err := template.ParseFS(pomXmlTemplate, "templates/pom_template.xml")
	if err != nil {
		fmt.Print("Failed to parse template")
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "pom_template.xml", xmlTemplate2)
	if err != nil {
		fmt.Print("Failed to execute template")
	}

	file, err := os.Create("output/pom.xml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, xmlTemplate2)
	if err != nil {
		fmt.Print("Failed to execute template, generate output file")
	}
}
