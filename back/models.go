package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
)

type Project struct {
	Id                  bson.ObjectId      `bson:"_id,omitempty" json:"id"`
	Name                string             `bson:"name" json:"name"`
	RepositoryType      string             `bson:"repository_type" json:"repository_type"`
	RepositoryUrl       string             `bson:"repository_url" json:"repository_url"`
	RepositoryToken     string             `bson:"repository_token" json:"repository_token,omitempty"`
	ExternalDomainNames []string           `bson:"external_domain_names" json:"external_domain_names"`
	UseHttps            bool               `bson:"use_https" json:"use_https"`
	TemplateId          string             `bson:"template_id" json:"template_id"`
	DockerImageVersion  string             `bson:"docker_image_version" json:"docker_image_version,omitempty"`
	RootFolder          string             `bson:"root_folder" json:"root_folder,omitempty"`
	ExposedPort         int                `bson:"exposed_port" json:"exposed_port,omitempty"`
	SystemDependencies  []string           `bson:"system_dependencies" json:"system_dependencies"`
	DependencyFiles     []string           `bson:"dependency_files" json:"dependency_files"`
	InstallScript       string             `bson:"install_script" json:"install_script,omitempty"`
	BuildScript         string             `bson:"build_script" json:"build_script,omitempty"`
	StartScript         string             `bson:"start_script" json:"start_script,omitempty"`
	StaticFolder        string             `bson:"static_folder" json:"static_folder,omitempty"`
	Databases           []*ProjectDatabase `bson:"databases" json:"databases"`
	Env                 map[string]string  `bson:"env" json:"env"`
	BelongsToId         string             `bson:"belongs_to" json:"-"`
	BelongsTo           ProjectUser        `bson:"-" json:"belongs_to"`
}

type ProjectDatabase struct {
	Type          string `bson:"type" json:"type"`
	Version       string `bson:"version" json:"version"`
	LinkDirectory string `bson:"-" json:"-"`
	EnvDatabase   string `bson:"env_db" json:"env_db"`
	EnvUsername   string `bson:"env_user" json:"env_user"`
	EnvPassword   string `bson:"env_password" json:"env_password"`
}

type ProjectUser struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Login     string        `json:"login"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Email     string        `json:"email"`
	Scope     string        `json:"scope"`
}

type ProjectFiles struct {
	Id            string            `json:"id"`
	ProjectId     string            `json:"project_id"`
	Dockerfiles   map[string]string `json:"dockerfiles"`
	DockerCompose string            `json:"docker_compose"`
}

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Id        string `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	MongoId   string `json:"mongo_id"`
	Scope     string `json:"scope"`
	jwt.StandardClaims
}