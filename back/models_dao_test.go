package main

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var u0 = ProjectUser{
	// Id not set
	OauthId:   0,
	Login:     "2019test",
	FirstName: "Testy",
	LastName:  "McTestface",
	Email:     "test@example.com",
	Scope:     "admin",
}

var d01 = ProjectDatabase{
	Type:        "mysql",
	Version:     "5.7",
	EnvDatabase: "test",
	EnvUsername: "root",
	EnvPassword: "password",
}

var d02 = ProjectDatabase{
	Type:        "postgres",
	Version:     "",
	EnvDatabase: "tp",
	EnvUsername: "rootp",
	EnvPassword: "passwordp",
}

var p0 = Project{
	// Id not set
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
	Databases:           []*ProjectDatabase{&d01, &d02},
	Env:                 map[string]string{"test1": "aaa", "test2": "bbb"},
	BelongsToId:         "5c509e17910f118485a1a7ba",
	BelongsTo:           &u0,
}

func getProject() Project {
	project := p0
	project.Id = bson.NewObjectId()
	user := u0
	user.Id = bson.NewObjectId()
	project.BelongsToId = user.Id.Hex()
	project.BelongsTo = &user
	return project
}

// DAO to access data from the database
var daoTest = TyphoonDAO{}

// Connect the DAO to the database if it is empty
func (td *TyphoonDAO) connectIfNeeded() {
	if daoTest.Database != "typhoon_test" {
		daoTest.Server = "mongodb://root:example@mongo:27017/"
		daoTest.Database = "typhoon_test"
		daoTest.Connect()
	}
}

// Clear all database. Use with caution
func (td *TyphoonDAO) clearDatabase() {
	db.C("projects").RemoveAll(bson.M{})
	db.C("users").RemoveAll(bson.M{})
}

// Tests on the functions using the projects collection
func TestProjectActions(t *testing.T) {
	daoTest.connectIfNeeded()
	daoTest.clearDatabase()

	var project Project
	var projects []Project
	var err error
	var found bool
	p1 := getProject()

	// Inserting user u1
	_, err = daoTest.InsertUser(*(p1.BelongsTo))
	if err != nil {
		t.Errorf("Error while inserting user: %s", err.Error())
	} else {
		t.Log("User inserted without error")
	}

	// Inserting project p1
	err = daoTest.InsertProject(p1)
	if err != nil {
		t.Errorf("Error while inserting project: %s", err.Error())
	} else {
		t.Log("Project inserted without error")
	}

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
		t.Errorf("Project in database messes BelongsToId: %s != %s", project.BelongsTo.Id.Hex(), project.BelongsToId)
	}
	if project.BelongsTo.Id != p1.BelongsTo.Id {
		t.Errorf("Project database BelongsTo.Id mismatch!: %s != %s", project.BelongsTo.Id, p1.BelongsTo.Id)
	}
	t.Log("Project p1 retrieved by id without fatal errors")

	// Get p1 back by name
	project, err = daoTest.FindProjectByName(p1.Name)
	if err != nil {
		t.Errorf("Error while searching for project by name: %s", err.Error())
	}
	if project.Id != p1.Id {
		t.Errorf("Project Id mismatch!: %s != %s", project.Id, p1.Id)
	}
	t.Log("Project p1 retrieved by name without fatal errors")

	// Check if p1 is in the list of all projects
	projects, err = daoTest.FindAllProjects()
	if err != nil {
		t.Errorf("Error while listing projects: %s", err.Error())
	}
	found = false
	for _, project := range projects {
		if project.Id == p1.Id {
			t.Log("Found p1 among list of projects")
			found = true
			break
		}
	}
	if !found {
		t.Error("p1 not found in list of projects")
	}

	// Check if p1 is in the list of projects of u1
	projects, err = daoTest.FindProjectsOfUser(p1.BelongsTo.Id.Hex())
	if err != nil {
		t.Errorf("Error while listing projects of u1: %s", err.Error())
	}
	found = false
	for _, project := range projects {
		if project.Id == p1.Id {
			t.Log("Found p1 among list of projects of u1")
			found = true
			break
		}
	}
	if !found {
		t.Error("p1 not found in list of projects of u1")
	}

	// Update p1 in database
	p1.Name = "unittestproject_2"
	err = daoTest.UpdateProject(p1)
	if err != nil {
		t.Errorf("Error while updating project p1: %s", err.Error())
	} else {
		t.Log("Updated p1 without errors")
	}

	// Get p1 back
	project, err = daoTest.FindProjectById(p1.Id.Hex())
	if err != nil {
		t.Errorf("Error while searching for project by id: %s", err.Error())
	}
	if project.Name != "unittestproject_2" {
		t.Errorf("Error: project.Name did not update? Value: %s", project.Name)
	}

	// Delete p1 in database
	err = daoTest.DeleteProject(p1.Id.Hex())
	if err != nil {
		t.Errorf("Error while deleting project p1: %s", err.Error())
	} else {
		t.Log("Deleted p1 without errors")
	}

	// Get p1 back
	project, err = daoTest.FindProjectById(p1.Id.Hex())
	if err == nil {
		t.Error("Error: p1 is still in db?")
	} else if err == mgo.ErrNotFound {
		t.Log("Deletion of p1 seems sussessful")
	} else {
		t.Errorf("Error when checking for project p1 absence: %s", err.Error())
	}
}

// Tests on the functions using the users collection
func TestUserActions(t *testing.T) {
	daoTest.connectIfNeeded()
	daoTest.clearDatabase()

	// var project Project
	// var projects []Project
	var user ProjectUser
	var users []ProjectUser
	var err error
	var found bool
	p1 := getProject()
	var u1 ProjectUser

	// Inserting user u1
	u1, err = daoTest.InsertUser(*(p1.BelongsTo))
	if err != nil {
		t.Errorf("Error while inserting user: %s", err.Error())
	} else {
		t.Log("User inserted without error")
	}

	// Inserting project p1
	err = daoTest.InsertProject(p1)
	if err != nil {
		t.Errorf("Error while inserting project: %s", err.Error())
	} else {
		t.Log("Project inserted without error")
	}

	// Get u1 back by id
	user, err = daoTest.FindUserById(u1.Id.Hex())
	if err != nil {
		t.Errorf("Error while searching for user by id: %s", err.Error())
	}
	if user.Id != u1.Id {
		t.Errorf("User Id mismatch!: %s != %s", user.Id, u1.Id)
	}
	if user.Login != u1.Login {
		t.Errorf("User Login mismatch!: %s != %s", user.Login, u1.Login)
	}
	t.Log("User u1 retrieved by id without fatal errors")

	// Get u1 back by login
	user, err = daoTest.FindUserByLogin(u1.Login)
	if err != nil {
		t.Errorf("Error while searching for user by login: %s", err.Error())
	}
	if user.Id != u1.Id {
		t.Errorf("User Id mismatch!: %s != %s", user.Id, u1.Id)
	}
	if user.Login != u1.Login {
		t.Errorf("User Login mismatch!: %s != %s", user.Login, u1.Login)
	}
	t.Log("User u1 retrieved by login without fatal errors")

	// Check if u1 is in the list of all users
	users, err = daoTest.FindAllUsers()
	if err != nil {
		t.Errorf("Error while listing users: %s", err.Error())
	}
	found = false
	for _, user := range users {
		if user.Id == u1.Id {
			t.Log("Found u1 among list of users")
			found = true
			break
		}
	}
	if !found {
		t.Error("u1 not found in list of users")
	}

	// Check if u1 is in the list of all admin users
	users, err = daoTest.FindAllAdminUsers()
	if err != nil {
		t.Errorf("Error while listing users: %s", err.Error())
	}
	found = false
	for _, user := range users {
		if user.Id == u1.Id {
			t.Log("Found u1 among list of admin users")
			found = true
			break
		}
	}
	if !found {
		t.Error("u1 not found in list of admin users")
	}

	// Update u1 in database
	u1.FirstName = "Testy_2"
	err = daoTest.UpdateUser(u1)
	if err != nil {
		t.Errorf("Error while updating user u1: %s", err.Error())
	} else {
		t.Log("Updated u1 without errors")
	}

	// Get u1 back
	user, err = daoTest.FindUserById(u1.Id.Hex())
	if err != nil {
		t.Errorf("Error while searching for project by id: %s", err.Error())
	}
	if user.FirstName != "Testy_2" {
		t.Errorf("Error: user.FirstName did not update? Value: %s", user.FirstName)
	}

	// Delete u1 in database
	err = daoTest.DeleteUser(u1.Id.Hex())
	if err != nil {
		t.Errorf("Error while deleting user u1: %s", err.Error())
	} else {
		t.Log("Deleted u1 without errors")
	}

	// Get u1 back
	user, err = daoTest.FindUserById(u1.Id.Hex())
	if err == nil {
		t.Error("Error: u1 is still in db?")
	} else if err == mgo.ErrNotFound {
		t.Log("Deletion of u1 seems sussessful")
	} else {
		t.Errorf("Error when checking for project u1 absence: %s", err.Error())
	}
}
