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
	TemplateFile     string `json:"template_file"`
	ImageSuffix      string `json:"image_name"`
	DockerfilePath   string `json:"dockerfile_path"`
	DockerfileExists bool   `json:"dockerfile_exists"`
}

type DockerData struct {
	Dockerfiles   []DockerfileData `json:"dockerfiles"`
	DockerCompose string           `json:"docker_compose"`
}

// Makes the mapping between the chosen project template and template files
var templateIdToFiles = map[string]DockerData{
	"node": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{TemplateFile: "dockerfiles/node.gotmpl", ImageSuffix: "-node"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"php": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{TemplateFile: "dockerfiles/php.gotmpl", ImageSuffix: "-php"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"python3": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{TemplateFile: "dockerfiles/python3.gotmpl", ImageSuffix: "-python3"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
	"react": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{TemplateFile: "dockerfiles/react.gotmpl", ImageSuffix: "-react"},
		},
		DockerCompose: "docker_composes/standalone.gotmpl",
	},
	"static": DockerData{
		Dockerfiles: []DockerfileData{
			DockerfileData{TemplateFile: "dockerfiles/static.gotmpl", ImageSuffix: "-static"},
		},
		DockerCompose: "docker_composes/standard.gotmpl",
	},
}

// From a project, a template and output info, writes and returns the filled template
func MakeStringAndFile(p interface{}, templateFile string, outputFile string) (string, error) {
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

	// Writing to file
	if outputFile != "" {
		log.Println("Will try to write " + templateFile + " filled in " + outputFile + "...")

		// Creates the target directory if it does not exist
		os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
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

// Get the code log path
func (p *Project) LogPath() string {
	return filepath.Join("/typhoon_logs", p.Id.Hex())
}

// Get the code clone paths
func (p *Project) ClonePath() string {
	return filepath.Join("/typhoon_sites", p.Id.Hex())
}

// Get the Dockerfile paths, and if they exist
func (p *Project) DockerfilePaths() ([]DockerfileData, error) {
	dd, ok := templateIdToFiles[p.TemplateId]
	if !ok {
		return nil, errors.New("Unknown template: " + p.TemplateId)
	}

	// Dockerfiles
	dockerfileData := make([]DockerfileData, 0)
	for _, dfd := range dd.Dockerfiles {
		outputDirectory := filepath.Join("/typhoon_dockerfile", p.Id.Hex()+dfd.ImageSuffix)
		outputFile := filepath.Join(outputDirectory, "Dockerfile")
		dfd.DockerfilePath = outputFile
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			dfd.DockerfileExists = false
		} else {
			dfd.DockerfileExists = true
		}
		dockerfileData = append(dockerfileData, dfd)
	}
	return dockerfileData, nil
}

// Get the docker-compose.yml directory, path, and if it exists
func (p *Project) DockerComposePaths() (string, string, bool) {
	dockerComposeFileDir := filepath.Join("/typhoon_docker_compose", p.Id.Hex())
	dockerComposeFilePath := filepath.Join(dockerComposeFileDir, "docker-compose.yml")
	if _, err := os.Stat(dockerComposeFilePath); os.IsNotExist(err) {
		return dockerComposeFileDir, dockerComposeFilePath, false // File is not already here
	}
	return dockerComposeFileDir, dockerComposeFilePath, true // File already exists
}
