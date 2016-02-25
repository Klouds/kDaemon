package database

import (
	"errors"
	r "github.com/dancannon/gorethink"
	"github.com/klouds/kDaemon/config"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
)

var (
	Session *r.Session
)

//Initializes supporting functions
func Init() {
	InitDB()
}

/* DATABASE FUNCTIONALITY */
// connect to the db
func InitDB() {

	logging.Log("Initializing Database connection.")

	rethinkdbhost, err := config.Config.GetString("default", "rethinkdb_host")
	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_host)")
	}

	rethinkdbport, err := config.Config.GetString("default", "rethinkdb_port")
	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_port)")
	}

	rethinkdbname, err := config.Config.GetString("default", "rethinkdb_dbname")

	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_dbname)")
	}

	session, err := r.Connect(r.ConnectOpts{
		Address:  rethinkdbhost + ":" + rethinkdbport,
		Database: rethinkdbname,
	})

	Session = session

}

//Node database functions

//Create a new node in the database
func CreateNode(n *models.Node) (string, error) {

	res, err := r.Table("nodes").
		Insert(n).
		RunWrite(Session)

	if err != nil {
		return "", err
	}

	containerid := ""

	keys := res.GeneratedKeys
	if len(keys) > 0 {
		containerid = keys[0]
	}
	return containerid, nil
}

//delete Node
func DeleteNode(id string) (bool, error) {

	err := r.Table("nodes").Get(id).Delete().Exec(Session)

	if err != nil {
		return false, err
	}

	return true, err

}

//Get node information
func GetNode(id string) (*models.Node, error) {

	res, err := r.Table("nodes").Get(id).Run(Session)

	if err != nil {
		return nil, err
	}

	var node models.Node

	err = res.One(&node)

	if err != nil {
		return nil, err
	}

	return &node, err
}

func GetNodes() ([]models.Node, error) {

	res, err := r.Table("nodes").Run(Session)

	if err != nil {
		return nil, err
	}

	var nodes []models.Node

	err = res.All(&nodes)

	if err != nil {
		return nil, err
	}

	if len(nodes) <= 0 {
		return nil, errors.New("NO NODES")
	}

	return nodes, err
}

//Update node
func UpdateNode(node *models.Node) (bool, error) {

	_, err := r.Table("nodes").
		Get(node.Id).
		Update(node).
		RunWrite(Session)

	if err != nil {
		return false, err
	}

	return true, err
}

//Application database functions

//Create a new application in the database
func CreateApplication(a *models.Application) (bool, error) {

	err := r.Table("applications").
		Insert(a).
		Exec(Session)
	if err != nil {
		return false, err
	}

	return true, err
}

//Get application information
func GetApplication(id string) (*models.Application, error) {

	res, err := r.Table("applications").Get(id).Run(Session)

	if err != nil {
		return nil, err
	}

	var app models.Application

	err = res.One(&app)

	if err != nil {
		return nil, err
	}

	return &app, err
}

func GetApplications() ([]models.Application, error) {
	//Returns a list of all applications

	res, err := r.Table("applications").Run(Session)

	if err != nil {
		return nil, err
	}

	var applications []models.Application

	err = res.All(&applications)

	if err != nil {
		return nil, err
	}

	if len(applications) <= 0 {
		return nil, errors.New("NO APPLICATIONS")
	}

	return applications, err
}

//delete application from database
func DeleteApplication(id string) (bool, error) {

	err := r.Table("applications").Get(id).Delete().Exec(Session)

	if err != nil {
		return false, err
	}

	return true, err

}

//Update Application
func UpdateApplication(app *models.Application) (bool, error) {

	_, err := r.Table("applications").
		Get(app.Id).
		Update(app).
		RunWrite(Session)

	if err != nil {
		return false, err
	}

	return true, err
}

//Container database Functions

//Create a new node in the database
func CreateContainer(c *models.Container) (*models.Container, bool, error) {
	//TODO: Check for auth

	resp, err := r.Table("containers").
		Insert(c).
		RunWrite(Session)
	if err != nil {
		return nil, false, err
	}

	if len(resp.GeneratedKeys) != 0 {
		c.Id = resp.GeneratedKeys[0]
	}

	return c, true, err
}

func UpdateContainer(cont *models.Container) (bool, error) {

	_, err := r.Table("containers").
		Get(cont.Id).
		Update(cont).
		RunWrite(Session)

	if err != nil {
		return false, err
	}

	return true, err
}

func GetContainerByName(name string) *models.Container {
	logging.Log("Getting container by name: ", name)
	//Look for a container with name
	var newcontainer models.Container

	resp, err := r.Table("containers").Filter(r.Row.Field("name").
		Eq(name)).Run(Session)

	if err != nil {
		return nil
	}

	err = resp.One(&newcontainer)

	if err != nil {
		return nil
	}

	return &newcontainer
}

//Get container information
func GetContainer(id string) (*models.Container, error) {

	res, err := r.Table("containers").Get(id).Run(Session)

	if err != nil {
		return nil, err
	}

	var container models.Container

	err = res.One(&container)

	if err != nil {
		return nil, err
	}

	return &container, err
}

func GetContainers() ([]models.Container, error) {
	//Returns a list of all applications

	res, err := r.Table("containers").Run(Session)

	if err != nil {
		return nil, err
	}

	var containers []models.Container

	err = res.All(&containers)

	if err != nil {
		return nil, err
	}

	if len(containers) <= 0 {
		return nil, errors.New("NO CONTAINERS")
	}

	return containers, err
}

func DeleteContainer(id string) error {

	err := r.Table("containers").Get(id).Delete().Exec(Session)

	return err

}
