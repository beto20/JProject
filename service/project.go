package service

type project struct {
	group       string
	artifact    string
	name        string
	description string
	packageName string
	javaVersion string
}

type ApplicationTemplate struct {
	Namespace    string
	BasePackages string
	Name         string
}

type projectService interface {
	generate(project project) (bool, error)
}

func NewProject() projectService {
	return &project{}
}

func (p *project) generate(project project) (bool, error) {

	return true, nil
}

// func buildApplication(path string) {
// 	applicationTmpl := ApplicationTemplate{
// 		Namespace:    "pe.com.mock",
// 		BasePackages: "pe.com",
// 		Name:         "MockApplication",
// 	}

// 	tmpl, err := template.ParseFS(applicationTemplate, "templates/main_application.tmpl")
// 	if err != nil {
// 		panic(err)
// 	}

// 	file, err := os.Create(path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = tmpl.Execute(file, applicationTmpl)
// 	if err != nil {
// 		panic(err)
// 	}
// }
