package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// From a project, will clone or pull the source
func GetSourceCode(p *Project) error {
	if p.RepositoryUrl == "" {
		return errors.New("RepositoryUrl not specified")
	}

	baseDir := "/typhoon_sites"
	clonePath := filepath.Join(baseDir, p.Id.Hex())

	log.Println("Will try to clone " + p.RepositoryUrl + " in " + clonePath + "...")

	os.RemoveAll(clonePath)
	os.MkdirAll(clonePath, os.ModePerm)

	cmdGit := exec.Command("git", "clone", "-q", "--depth", "1", "--", p.RepositoryUrl, clonePath)
	cmdGit.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmdGit.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// From a project, will write all the templates in files
func WriteFromTemplates(p *Project) (map[string]string, error) {
	dd, ok := templateIdToFiles[p.TemplateId]
	if !ok {
		return nil, errors.New("Unknown template: " + p.TemplateId)
	}

	outputDirectory := ""

	ps, err := json.Marshal(p)
	results := map[string]string{"project": string(ps)}

	// Dockerfiles
	for i, dfd := range dd.Dockerfiles {
		outputDirectory = filepath.Join("/typhoon_dockerfile", p.Id.Hex()+dfd.ImageName)
		res, err := MakeStringAndFile(p, dfd.TemplateFile, outputDirectory, "Dockerfile")
		results[fmt.Sprintf("dockerfile_%d", i)] = res
		if err != nil {
			results[fmt.Sprintf("error_dockerfile_%d", i)] = err.Error()
		}
	}

	// docker-compose file
	outputDirectory = filepath.Join("/typhoon_docker_compose", p.Id.Hex())
	res, err := MakeStringAndFile(p, dd.DockerCompose, outputDirectory, "docker-compose.yml")
	results["docker_compose"] = res
	if err != nil {
		results["error_docker_compose"] = err.Error()
	}

	return results, nil
}
