package main

import (
	"testing"
)

// DAO to access data from the database
var daoTest = TyphoonDAO{}

func (td *TyphoonDAO) connectIfNeeded() {
	if daoTest.Database != "typhoon_test" {
		// Create the DAO object and connect it to the mongo server
		daoTest.Server = "mongodb://root:example@mongo:27017/"
		daoTest.Database = "typhoon_test"
		daoTest.Connect()
	}
}

var u1 = ProjectUser{
	Id:        "5c509e17910f118485a1a7ba",
	Login:     "2015bernarda",
	FirstName: "Aymeric",
	LastName:  "Bernard",
	Email:     "nope@nope.fr",
	Scope:     "admin",
}

var d1 = ProjectDatabase{
	Type:        "mysql",
	Version:     "5.7",
	EnvDatabase: "test",
	EnvUsername: "root",
	EnvPassword: "password",
}

var d2 = ProjectDatabase{
	Type:        "postgres",
	Version:     "",
	EnvDatabase: "tp",
	EnvUsername: "rootp",
	EnvPassword: "passwordp",
}

var p1 = Project{
	Id:                  "123456",
	Name:                "unittestproject",
	RepositoryType:      "github",
	RepositoryUrl:       "https://github.com/typhoon-docker/example-flask",
	RepositoryToken:     "",
	ExternalDomainNames: []string{"unittestproject.fr", "unittestproject.com"},
	UseHttps:            true,
	TemplateId:          "python3",
	DockerImageVersion:  "3.7",
	RootFolder:          "",
	ExposedPort:         8042,
	SystemDependencies:  []string{"git", "screen"},
	DependencyFiles:     []string{"requirements.txt"},
	InstallScript:       "echo installing",
	BuildScript:         "echo building",
	StartScript:         "python3 flaskserver.py",
	StaticFolder:        "",
	Databases:           []*ProjectDatabase{&d1, &d2},
	Env:                 map[string]string{"test1": "aaa", "test2": "bbb"},
	BelongsToId:         "5c509e17910f118485a1a7ba",
	BelongsTo:           &u1,
}

// func (m *TyphoonDAO) FindAllProjects() ([]Project, error)
// func (m *TyphoonDAO) FindProjectsOfUser(uMongoId string) ([]Project, error)
// func (m *TyphoonDAO) FindProjectById(id string) (Project, error)
// func (m *TyphoonDAO) FindProjectByName(name string) (Project, error)
// func (m *TyphoonDAO) InsertProject(project Project) error
// func (m *TyphoonDAO) DeleteProject(id string) error
// func (m *TyphoonDAO) UpdateProject(project Project) error
// func (m *TyphoonDAO) FindUserById(id string) (ProjectUser, error)
// func (m *TyphoonDAO) FindUserByLogin(login string) (ProjectUser, error)
// func (m *TyphoonDAO) InsertUser(user ProjectUser) (ProjectUser, error)
// func (m *TyphoonDAO) FindAllUsers() ([]ProjectUser, error)
// func (m *TyphoonDAO) FindAllAdminUsers() ([]ProjectUser, error)
// func (m *TyphoonDAO) UpdateUser(user ProjectUser) error
// func (m *TyphoonDAO) DeleteUser(id string) error

func TestProjectActions(t *testing.T) {
	var project Project
	var projects []Project
	var err error

	// Inserting project p1
	err = daoTest.InsertProject(p1)
	if err != nil {
		t.Errorf("Error while inserting project: %s", err.Error())
	}
	t.Log("Project inserted without error")

	// Get p1 back
	project, err = daoTest.FindProjectById(p1.Id.Hex())
	if err != nil {
		t.Errorf("Error while searching for project by id: %s", err.Error())
	}
	if project.Id != p1.Id {
		t.Errorf("Project Id mismatch!: %s != %s", project.Id, p1.Id)
	}
	if project.Name != p1.Name {
		t.Errorf("Project Name mismatch!: %s != %s", project.Name, p1.Name)
	}
	if project.Databases[0].Type != p1.Databases[0].Type {
		t.Errorf("Project database Type mismatch!: %s != %s", project.Databases[0].Type, p1.Databases[0].Type)
	}
	if project.BelongsToId != project.BelongsTo.Id.Hex() {
		t.Errorf("Project in database messes BelongsToId: %s != %s", project.BelongsToId, project.BelongsTo.Id.Hex())
	}
	if project.BelongsTo.Id != p1.BelongsTo.Id {
		t.Errorf("Project database BelongsTo.Id mismatch!: %s != %s", project.BelongsTo.Id, p1.BelongsTo.Id)
	}

	// Get p1 back by name
	project, err = daoTest.FindProjectByName(p1.Name)
	if err != nil {
		t.Errorf("Error while searching for project by name: %s", err.Error())
	}
	if project.Id != p1.Id {
		t.Errorf("Project Id mismatch!: %s != %s", project.Id, p1.Id)
	}

	// Check if p1 is in the list of projects
	projects, err = daoTest.FindAllProjects()
	if err != nil {
		t.Errorf("Error while listing projects: %s", err.Error())
	}
	if project.Id != p1.Id {
		t.Errorf("Project Id mismatch!: %s != %s", project.Id, p1.Id)
	}
}
