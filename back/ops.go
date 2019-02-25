package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// The project's host. Used by the docker-compose templates for urls
func (p *Project) Host() string {
	return os.Getenv("HOST")
}

// Find the containers id used by the project
func FindContainerID(containerImage string) (string, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		log.Println("FindContainerID NewClientWithOpts error: " + err.Error())
		return "", errors.New("FindContainerID NewClientWithOpts error: " + err.Error())
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Println("FindContainerID ContainerList error: " + err.Error())
		return "", errors.New("FindContainerID ContainerList error: " + err.Error())
	}

	for _, container := range containers {
		// fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		if containerImage == container.Image {
			return container.ID, nil
		}
	}
	return "", errors.New("Could not find container")
}

// Get containers info for the project
func GetContainerStatus(p *Project) ([]types.Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := client.NewClientWithOpts(client.WithVersion("1.39"))

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	projectContainers := make([]types.Container, 0)

	for _, container := range containers {
		// fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		if p.Name == container.Image {
			projectContainers = append(projectContainers, container)
		}
	}

	return projectContainers, nil
}

// GetContainerStats such as mem usage
func GetContainerStats(p *Project) (string, error) {
	containerID, err := FindContainerID(p.Name)
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := client.NewClientWithOpts(client.WithVersion("1.39"))

	stats, err := client.ContainerStats(ctx, containerID, false)

	if err != nil {
		log.Println(err)
	}

	log.Println(stats.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(stats.Body)
	statsReponse := buf.String()

	return statsReponse, nil

}

// Get the logs of the project
func GetLogsByName(p *Project, lines int) (string, error) {
	// Get the last timestamp
	timestamp := readLastLineTimestamp(p.LogPath())

	logs_file, err := os.OpenFile(p.LogPath(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer logs_file.Close()
	if err != nil {
		log.Println("GetLogsByName OpenFile: " + err.Error())
		return "", errors.New("GetLogsByName OpenFile: " + err.Error())
	}

	containerID, err := FindContainerID(p.Name)
	if err != nil {
		log.Println("GetLogsByName FindContainerID:" + err.Error())
		return "", errors.New("GetLogsByName FindContainerID: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := client.NewClientWithOpts(client.WithVersion("1.39"))
	reader, err := client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true, Since: timestamp})
	if err != nil {
		log.Println("GetLogsByName ContainerLogs: " + err.Error())
		return "", errors.New("GetLogsByName ContainerLogs: " + err.Error())
	}

	stdout := bytes.NewBuffer(make([]byte, 0))
	stdcopy.StdCopy(stdout, stdout, reader)

	if _, err = logs_file.WriteString(stdout.String()); err != nil {
		log.Println("GetLogsByName WriteString: " + err.Error())
		return "", errors.New("GetLogsByName WriteString: " + err.Error())
	}

	return ReadLastLines(p.LogPath(), lines), nil
}

// Clean the logs of the project
func CleanLogs(p *Project) error {
	os.Remove(p.LogPath())
	return nil
}

// From a project, will clone or pull the source
func GetSourceCode(p *Project) (map[string]string, error) {
	outputs := map[string]string{}

	if p.RepositoryUrl == "" {
		return outputs, errors.New("RepositoryUrl not specified")
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
	if p.RepositoryToken != "" {
		repoUrl = strings.Replace(repoUrl, "https://", "https://"+p.RepositoryToken+"@", 1)
	}
	branch := p.Branch
	if branch == "" {
		branch = "master"
	}
	cmd := exec.Command("git", "clone", "-b", branch, "--single-branch", "-q", "--depth", "1", "--", repoUrl, clonePath)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	outputs["git_clone"] = string(out)
	if err != nil {
		log.Println(cmd)
		log.Println("Could not clone: " + err.Error())
		outputs["error_git_clone"] = err.Error()
		return outputs, err
	} else {
		outputs["error_git_clone"] = ""
	}
	return outputs, nil
}

// From a project, will delete the clone directory
func CleanSourceCode(p *Project) error {
	os.RemoveAll(p.ClonePath())
	return nil
}

// From a project, will write all the templates in files
func FillTemplates(p *Project, write bool) (map[string]string, error) {
	dd, ok := templateIdToFiles[p.TemplateId]
	if !ok {
		return nil, errors.New("Unknown template: " + p.TemplateId)
	}

	// Results will hold template results, errors, and the project JSON itself
	ps, err := json.Marshal(p)
	results := map[string]string{"project": string(ps)}
	outputFile := ""

	// Dockerfiles
	dockerfileDataA, _ := p.DockerfilePaths()

	for i, dfd := range dockerfileDataA {
		if write {
			outputFile = dfd.DockerfilePath
		}
		res, err := MakeStringAndFile(p, dfd.TemplateFile, outputFile)
		results[fmt.Sprintf("dockerfile_%d", i)] = res
		if err != nil {
			results[fmt.Sprintf("error_dockerfile_%d", i)] = err.Error()
		}
	}

	// docker-compose file
	if write {
		_, outputFile, _ = p.DockerComposePaths()
	}
	res, err := MakeStringAndFile(p, dd.DockerCompose, outputFile)
	results["docker_compose"] = res
	if err != nil {
		log.Println("Could not create from template: " + err.Error())
		results["error_docker_compose"] = err.Error()
	}

	return results, nil
}

// From a project, will clean all docker files
func CleanDockerFiles(p *Project) error {
	// Dockerfiles
	dockerfileDataA, _ := p.DockerfilePaths()
	for _, dfd := range dockerfileDataA {
		os.RemoveAll(filepath.Dir(dfd.DockerfilePath))
	}
	// docker-compose
	dfDir, _, _ := p.DockerComposePaths()
	os.RemoveAll(dfDir)

	return nil
}

// From a project, will write all the templates in files
func BuildImages(p *Project) (map[string]string, error) {
	outputs := map[string]string{}

	// Gather Dockerfile data
	dockerfileDataA, err := p.DockerfilePaths()
	if err != nil {
		return outputs, err
	}

	// Build from Dockerfiles
	for i, dfd := range dockerfileDataA {
		// Location of the code, from where the Dockerfile is based
		context := filepath.Join(p.ClonePath(), p.RootFolder)

		// Run command to build. Uses the host's /var/run/docker.sock to build image into host
		log.Println("Will try to build " + p.Name + " from " + dfd.DockerfilePath + "...")
		cmd := exec.Command("docker", "build", "-t", p.Name, "-f", dfd.DockerfilePath, context)
		out, err := cmd.CombinedOutput()
		outputs[fmt.Sprintf("dockerfile_%d", i)] = string(out)
		if err != nil {
			log.Println("Could not build image: " + err.Error())
			outputs[fmt.Sprintf("error_dockerfile_%d", i)] = err.Error()
			return outputs, err
		} else {
			outputs[fmt.Sprintf("error_dockerfile_%d", i)] = ""
		}
	}
	return outputs, nil
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
		return errors.New("Could not run docker-compose down: " + err.Error())
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

// From a project, will delete the persistent directory
func CleanVolume(p *Project) error {
	os.RemoveAll(p.VolumePath())
	return nil
}
