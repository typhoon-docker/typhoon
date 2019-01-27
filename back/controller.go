package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type person struct {
	Name  string
	Phone string
}

func test() {
	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	fmt.Println("test")
	session, err := mgo.Dial("mongodb://root:example@mongo:27017/")
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("projects")
	err = c.Insert(
		&Project{
			Name:           "sade",
			RepositoryType: "github",
		},
		&Project{
			Name:           "delthasss",
			RepositoryType: "github",
		})
	if err != nil {
		log.Fatal(err)
	}

	result := Project{}
	err = c.Find(bson.M{"name": "sade"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("project:", result.RepositoryType)
	return
}
