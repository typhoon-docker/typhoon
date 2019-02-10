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
	// Where to clone the code
	baseDir := "/typhoon_sites"
	clonePath := filepath.Join(baseDir, p.Id.Hex())

	log.Println("Will try to clone " + p.RepositoryUrl + " in " + clonePath + "...")

	// Clean the target directory (maybe pull if already cloned?)
	os.RemoveAll(clonePath)
	os.MkdirAll(clonePath, os.ModePerm)

	// Run the clone command
	cmdGit := exec.Command("git", "clone", "-q", "--depth", "1", "--", p.RepositoryUrl+".git", clonePath)
	cmdGit.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmdGit.Run(); err != nil {
		log.Println("Could not clone: " + err.Error())
		return err
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

	// results will hold template results, errors, and the project JSON itself
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
		log.Println("Could not create from template: " + err.Error())
		results["error_docker_compose"] = err.Error()
	}

	return results, nil
}

// From a project, will write all the templates in files
func BuildImages(p *Project) error {
	dd, ok := templateIdToFiles[p.TemplateId]
	if !ok {
		return errors.New("Unknown template: " + p.TemplateId)
	}

	// Build from Dockerfiles
	for _, dfd := range dd.Dockerfiles {
		// Location of the Dockerfile to build
		dockerfileDirectory := filepath.Join("/typhoon_dockerfile", p.Id.Hex()+dfd.ImageName)
		fileName := filepath.Join(dockerfileDirectory, "Dockerfile")

		// Location of the code, from where the Dockerfile is based
		context := filepath.Join("/typhoon_sites", p.Id.Hex(), p.RootFolder)

		// Run command to build. Uses the host's /var/run/docker.sock to build image into host
		log.Println("Will try to build " + p.Name + " from " + fileName + "...")
		cmd := exec.Command("docker", "build", "-t", p.Name, "-f", fileName, context)
		if err := cmd.Run(); err != nil {
			log.Println("Could build image: " + err.Error())
			return err
		}
	}
	return nil
}

// From a project, will write all the templates in files
func DockerUp(p *Project) error {
	dockerComposeFileDir := filepath.Join("/typhoon_docker_compose", p.Id.Hex())
	dockerComposeFilePath := filepath.Join(dockerComposeFileDir, "docker-compose.yml")
	if _, err := os.Stat(dockerComposeFilePath); os.IsNotExist(err) {
		return errors.New("docker-composse.yml file not found in: " + dockerComposeFilePath)
	}

	// Run command to up
	log.Println("Will try to up from " + dockerComposeFilePath + "...")
	cmd := exec.Command("docker-compose", "up", "-d") // -d ?
	cmd.Dir = dockerComposeFileDir
	if err := cmd.Run(); err != nil {
		return errors.New("Could not run docker-compose up: " + err.Error())
	}
	return nil
}

// From a project, will write all the templates in files
func DockerDown(p *Project) error {
	dockerComposeFileDir := filepath.Join("/typhoon_docker_compose", p.Id.Hex())
	dockerComposeFilePath := filepath.Join(dockerComposeFileDir, "docker-compose.yml")
	if _, err := os.Stat(dockerComposeFilePath); os.IsNotExist(err) {
		return errors.New("docker-composse.yml file not found in: " + dockerComposeFilePath)
	}

	// Run command to down
	log.Println("Will try to down from " + dockerComposeFilePath + "...")
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = dockerComposeFileDir
	if err := cmd.Run(); err != nil {
		return errors.New("Could run docker-compose down: " + err.Error())
	}
	return nil
}
