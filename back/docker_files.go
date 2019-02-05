package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// // TEMP: Example
// var u1 = ProjectUser{
// 	Id:        "5c509e17910f118485a1a7ba",
// 	Login:     "2015bernarda",
// 	FirstName: "aymeric",
// 	LastName:  "bernard",
// 	Email:     "nope@nope.fr",
// 	Scope:     "admin",
// }

// // TEMP: Example
// var d1 = ProjectDatabase{
// 	Type:        "mysql",
// 	Version:     "5.7",
// 	EnvDatabase: "test",
// 	EnvUsername: "root",
// 	EnvPassword: "password",
// }

// // TEMP: Example
// var d2 = ProjectDatabase{
// 	Type:        "postgres",
// 	Version:     "",
// 	EnvDatabase: "tp",
// 	EnvUsername: "rootp",
// 	EnvPassword: "passwordp",
// }

// // TEMP: Example
// var p1 = Project{
// 	Id:                  "123456",
// 	Name:                "goProject1",
// 	RepositoryType:      "github",
// 	RepositoryUrl:       "https://github.com/typhoon-docker/example-flask",
// 	RepositoryToken:     "",
// 	ExternalDomainNames: []string{"myflask.fr", "cake.fr"},
// 	UseHttps:            true,
// 	TemplateId:          "python3",
// 	DockerImageVersion:  "3.7",
// 	RootFolder:          "",
// 	ExposedPort:         8042,
// 	SystemDependencies:  []string{"git", "screen"},
// 	DependencyFiles:     []string{"requirements.txt", "nope.txt"},
// 	InstallScript:       "echo installing",
// 	BuildScript:         "echo building",
// 	StartScript:         "python3 flaskserver.py",
// 	StaticFolder:        "",
// 	Databases:           []*ProjectDatabase{&d1, &d2},
// 	Env:                 map[string]string{"test1": "lol", "test2": "mdr"},
// 	BelongsToId:         "5c509e17910f118485a1a7ba",
// 	BelongsTo:           u1,
// }

type DockerfileData struct {
	TemplateFile string `json:"template_file"`
	ImageName    string `json:"image_name"`
}

type DockerData struct {
	Dockerfiles   []DockerfileData `json:"dockerfiles"`
	DockerCompose string           `json:"docker_compose"`
}

// Makes the mapping between the chosen project template and template files
var templateIdToFiles = map[string]DockerData{
	"node": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{"dockerfiles/node.gotmpl", "-node"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"php": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{"dockerfiles/php.gotmpl", "-php"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"python3": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{"dockerfiles/python3.gotmpl", "-python3"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"react": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{"dockerfiles/react.gotmpl", "-react"},
		},
		DockerCompose: "docker_composes/standalone.gotmpl",
	},
	"static": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{"dockerfiles/static.gotmpl", "-static"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
}

// From a project, a template and output info, writes and returns the filled template
func MakeStringAndFile(p interface{}, templateFile string, outputDirectory string, fileName string) (string, error) {
	outputFile := filepath.Join(outputDirectory, fileName)
	log.Println("Will try to write " + templateFile + " filled in " + outputFile + "...")

	// Get the template, make the result
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", errors.New("MakeStringAndFile ParseFiles: " + err.Error())
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, p); err != nil {
		return "", errors.New("MakeStringAndFile Execute: " + err.Error())
	}
	result := tpl.String()
	// log.Print("Template result:\n", result) // TEMP

	// Writing to file
	if fileName != "" {
		os.MkdirAll(outputDirectory, os.ModePerm)
		f, err := os.Create(outputFile)
		if err != nil {
			return "", errors.New("MakeStringAndFile Create: " + err.Error())
		}
		if _, err := tpl.WriteTo(f); err != nil {
			return "", errors.New("MakeStringAndFile WriteTo: " + err.Error())
		} else {
			log.Println("Wrote in " + outputFile)
		}
		f.Close()
	}
	return result, nil
}
