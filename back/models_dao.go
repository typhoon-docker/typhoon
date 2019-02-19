package main

import (
	"errors"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TyphoonDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

// Establish a connection to database
func (d *TyphoonDAO) Connect() {
	session, err := mgo.Dial(d.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(d.Database)
	log.Println("DAO: mongo db connected")
}

//////////////////////////////
////////// PROJECTS //////////
//////////////////////////////

// Find list of project
func (m *TyphoonDAO) FindAllProjects() ([]Project, error) {
	projects := make([]Project, 0)
	err := db.C("projects").Find(bson.M{}).All(&projects)
	for i := range projects {
		userId := projects[i].BelongsToId
		user, _ := m.FindUserById(userId)
		projects[i].BelongsTo = &user
	}
	return projects, err
}

// Find list of project for given user id
func (m *TyphoonDAO) FindProjectsOfUser(uMongoId string) ([]Project, error) {
	projects := make([]Project, 0)
	err := db.C("projects").Find(bson.M{"belongs_to": uMongoId}).All(&projects)
	user, _ := m.FindUserById(uMongoId)
	for i := range projects {
		projects[i].BelongsTo = &user
	}
	return projects, err
}

// Find a project by its id
func (m *TyphoonDAO) FindProjectById(id string) (Project, error) {
	var project Project
	err1 := db.C("projects").FindId(bson.ObjectIdHex(id)).One(&project)
	if err1 != nil {
		return project, err1
	}
	user, _ := m.FindUserById(project.BelongsToId)
	project.BelongsTo = &user
	return project, nil
}

// Find a project by its name
func (m *TyphoonDAO) FindProjectByName(name string) (Project, error) {
	var project Project
	err1 := db.C("projects").Find(bson.M{"name": name}).One(&project)
	if err1 != nil {
		return project, err1
	}
	user, _ := m.FindUserById(project.BelongsToId)
	project.BelongsTo = &user
	return project, nil
}

// Find a project by its repository URL
func (m *TyphoonDAO) FindProjectsByUrl(repositoryUrl string) ([]Project, error) {
	projects := make([]Project, 0)
	err := db.C("projects").Find(bson.M{"repository_url": repositoryUrl}).All(&projects)
	for i := range projects {
		userId := projects[i].BelongsToId
		user, _ := m.FindUserById(userId)
		projects[i].BelongsTo = &user
	}
	return projects, err
}

// Insert a project into database
func (m *TyphoonDAO) InsertProject(project Project) error {
	err := db.C("projects").Insert(&project)
	return err
}

// Delete an existing project
func (m *TyphoonDAO) DeleteProject(id string) error {
	m.DeleteLogs(id)
	err := db.C("projects").RemoveId(bson.ObjectIdHex(id))
	return err
}

// Update an existing project
func (m *TyphoonDAO) UpdateProject(project Project) error {
	err := db.C("projects").UpdateId(project.Id, &project)
	return err
}

///////////////////////////
////////// USERS //////////
///////////////////////////

// Find a user by its id
func (m *TyphoonDAO) FindUserById(id string) (ProjectUser, error) {
	var user ProjectUser
	if id == "" {
		return user, errors.New("Trying to find user by id with empty id!")
	}
	err := db.C("users").FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Find a user by its login
func (m *TyphoonDAO) FindUserByLogin(login string) (ProjectUser, error) {
	var user ProjectUser
	err := db.C("users").Find(bson.M{"login": login}).One(&user)
	return user, err
}

// Insert a user in the database
func (m *TyphoonDAO) InsertUser(user ProjectUser) (ProjectUser, error) {
	if user.Id.Hex() == "" {
		user.Id = bson.NewObjectId()
	}
	err := db.C("users").Insert(&user)
	return user, err
}

// Find list of users
func (m *TyphoonDAO) FindAllUsers() ([]ProjectUser, error) {
	users := make([]ProjectUser, 0)
	err := db.C("users").Find(bson.M{}).All(&users)
	return users, err
}

// Find list of admin users
func (m *TyphoonDAO) FindAllAdminUsers() ([]ProjectUser, error) {
	users := make([]ProjectUser, 0)
	err := db.C("users").Find(bson.M{"scope": "admin"}).All(&users)
	return users, err
}

// Update an existing user
func (m *TyphoonDAO) UpdateUser(user ProjectUser) error {
	err := db.C("users").UpdateId(user.Id, &user)
	return err
}

// Delete user
func (m *TyphoonDAO) DeleteUser(id string) error {
	err := db.C("users").RemoveId(bson.ObjectIdHex(id))
	return err
}

//////////////////////////////////
////////// PROJECT LOGS //////////
//////////////////////////////////

// Find project logs by its id
func (m *TyphoonDAO) FindLogsByProjectId(id string) (ProjectLogs, error) {
	var logs ProjectLogs
	if id == "" {
		return logs, errors.New("Trying to find logs by id with empty id!")
	}
	err := db.C("logs").FindId(bson.ObjectIdHex(id)).One(&logs)
	return logs, err
}

// Insert project logs in the database
func (m *TyphoonDAO) InsertLogs(logs ProjectLogs) (ProjectLogs, error) {
	err := db.C("logs").Insert(&logs)
	return logs, err
}

// Update existing logs, adding new keys and overriding existing ones
func (m *TyphoonDAO) UpdateLogs(logs ProjectLogs) (ProjectLogs, error) {
	oldLogs, err := m.FindLogsByProjectId(logs.Id.Hex())
	if err == mgo.ErrNotFound {
		oldLogs = ProjectLogs{Id: logs.Id, Logs: map[string]string{}}
		m.InsertLogs(oldLogs)
	} else if err != nil {
		return oldLogs, err
	}
	for nlKey, nlValue := range logs.Logs {
		oldLogs.Logs[nlKey] = nlValue
	}
	err = db.C("logs").UpdateId(logs.Id, &oldLogs)
	return oldLogs, err
}

// Update existing logs, adding new keys and overriding existing ones
func (m *TyphoonDAO) UpdateLogsById(id string, content map[string]string) (ProjectLogs, error) {
	return m.UpdateLogs(ProjectLogs{Id: bson.ObjectIdHex(id), Logs: content})
}

// Delete logs
func (m *TyphoonDAO) DeleteLogs(id string) error {
	err := db.C("logs").RemoveId(bson.ObjectIdHex(id))
	return err
}
