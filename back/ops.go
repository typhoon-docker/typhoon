package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// From a project, will clone or pull the source
func GetSourceCode(p *Project) error {
	if p.RepositoryUrl == "" {
		return errors.New("RepositoryUrl not specified")
	}
	// Where to clone the code
	clonePath := p.ClonePath()

	log.Println("Will try to clone " + p.RepositoryUrl + " in " + clonePath + "...")

	// Clean the target directory (maybe pull if already cloned?)
	os.RemoveAll(clonePath)
	os.MkdirAll(clonePath, os.ModePerm)

	// Run the clone command
	repoUrl := p.RepositoryUrl
	if !strings.HasSuffix(repoUrl, ".git") {
		repoUrl = repoUrl + ".git"
	}
	cmdGit := exec.Command("git", "clone", "-q", "--depth", "1", "--", repoUrl, clonePath)
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

	// Results will hold template results, errors, and the project JSON itself
	ps, err := json.Marshal(p)
	results := map[string]string{"project": string(ps)}

	// Dockerfiles
	dockerfileDataA, _ := p.DockerfilePaths()

	for i, dfd := range dockerfileDataA {
		res, err := MakeStringAndFile(p, dfd.TemplateFile, dfd.DockerfilePath)
		results[fmt.Sprintf("dockerfile_%d", i)] = res
		if err != nil {
			results[fmt.Sprintf("error_dockerfile_%d", i)] = err.Error()
		}
	}

	// docker-compose file
	_, outputFile, _ := p.DockerComposePaths()
	res, err := MakeStringAndFile(p, dd.DockerCompose, outputFile)
	results["docker_compose"] = res
	if err != nil {
		log.Println("Could not create from template: " + err.Error())
		results["error_docker_compose"] = err.Error()
	}

	return results, nil
}

// From a project, will write all the templates in files
func BuildImages(p *Project) error {
	dockerfileDataA, err := p.DockerfilePaths()
	if err != nil {
		return err
	}

	// Build from Dockerfiles
	for _, dfd := range dockerfileDataA {
		// Location of the code, from where the Dockerfile is based
		context := filepath.Join(p.ClonePath(), p.RootFolder)

		// Run command to build. Uses the host's /var/run/docker.sock to build image into host
		log.Println("Will try to build " + p.Name + " from " + dfd.DockerfilePath + "...")
		cmd := exec.Command("docker", "build", "-t", p.Name, "-f", dfd.DockerfilePath, context)
		if err := cmd.Run(); err != nil {
			log.Println("Could build image: " + err.Error())
			return err
		}
	}
	return nil
}

// From a project, will run docker-compose up
func DockerUp(p *Project) error {
	dockerComposeFileDir, dockerComposeFilePath, here := p.DockerComposePaths()
	if !here {
		return errors.New("docker-compose.yml file not found in: " + dockerComposeFilePath)
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

// From a project, will run docker-compose down
func DockerDown(p *Project) error {
	dockerComposeFileDir, dockerComposeFilePath, here := p.DockerComposePaths()
	if !here {
		return errors.New("docker-compose.yml file not found in: " + dockerComposeFilePath)
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

// From a project, will run docker-compose ps
func DockerStatus(p *Project) (string, error) {
	dockerComposeFileDir, dockerComposeFilePath, here := p.DockerComposePaths()
	if !here {
		return "", errors.New("docker-compose.yml file not found in: " + dockerComposeFilePath)
	}

	// Run command to down
	log.Println("Will try to get status from " + dockerComposeFilePath + "...")
	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = dockerComposeFileDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New("Could run docker-compose ps: " + err.Error())
	}
	return string(out), nil
}
