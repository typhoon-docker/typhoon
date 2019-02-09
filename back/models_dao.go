package main

import (
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
	log.Println("DAO: db connected")
}

// Find list of project
func (m *TyphoonDAO) FindAllProjects() ([]Project, error) {
	var projects []Project
	err := db.C("projects").Find(bson.M{}).All(&projects)
	return projects, err
}

// Find list of project for given user id
func (m *TyphoonDAO) FindProjectsOfUser(uMongoId string) ([]Project, error) {
	var projects []Project
	err := db.C("projects").Find(bson.M{"belongs_to": uMongoId}).All(&projects)
	return projects, err
}

// Find a project by its id
func (m *TyphoonDAO) FindProjectById(id string) (Project, error) {
	var project Project
	err := db.C("projects").FindId(bson.ObjectIdHex(id)).One(&project)
	return project, err
}

// Find a project by its id
func (m *TyphoonDAO) FindProjectByName(name string) (Project, error) {
	var project Project
	err := db.C("projects").Find(bson.M{"name": name}).One(&project)
	return project, err
}

// Insert a project into database
func (m *TyphoonDAO) InsertProject(project Project) error {
	err := db.C("projects").Insert(&project)
	return err
}

// Delete an existing project
func (m *TyphoonDAO) DeleteProject(id string) error {
	err := db.C("projects").RemoveId(bson.ObjectIdHex(id))
	return err
}

// Update an existing project
func (m *TyphoonDAO) UpdateProject(project Project) error {
	err := db.C("projects").UpdateId(project.Id, &project)
	return err
}

// Find a user by its id
func (m *TyphoonDAO) FindUserById(id string) (ProjectUser, error) {
	var user ProjectUser
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
	user.Id = bson.NewObjectId()
	err := db.C("users").Insert(&user)
	return user, err
}

/////// TEMP ? ///////
// Find list of users
func (m *TyphoonDAO) FindAllUsers() ([]ProjectUser, error) {
	var users []ProjectUser
	err := db.C("users").Find(bson.M{}).All(&users)
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
