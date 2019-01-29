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
func (m *TyphoonDAO) FindProjectsOfUser(uid string) ([]Project, error) {
	var projects []Project
	err := db.C("projects").Find(bson.M{"belongs_to": uid}).All(&projects)
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
func (m *TyphoonDAO) DeleteProject(project Project) error {
	err := db.C("projects").RemoveId(project.Id)
	return err
}

// Update an existing project
func (m *TyphoonDAO) UpdateProject(project Project) error {
	err := db.C("projects").UpdateId(project.Id, &project)
	return err
}
