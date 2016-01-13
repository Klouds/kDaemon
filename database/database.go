package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/superordinate/kDaemon/config"
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/models"
)

var (
	db *gorm.DB
)

//Initializes supporting functions
func Init() {
	InitDB()
}

/* DATABASE FUNCTIONALITY */
// connect to the db
func InitDB() {

	logging.Log("Initializing Database connection.")

	mysqlhost, err := config.Config.GetString("default", "mysql_host")
	if err != nil {
		logging.Log("Problem with config file! (mysql_host)")
	}
	mysqlport, err := config.Config.GetString("default", "mysql_port")
	if err != nil {
		logging.Log("Problem with config file! (mysql_port)")
	}
	mysqluser, err := config.Config.GetString("default", "mysql_user")
	if err != nil {
		logging.Log("Problem with config file! (mysql_user)")
	}
	mysqlpass, err := config.Config.GetString("default", "mysql_password")
	if err != nil {
		logging.Log("Problem with config file! (mysql_port)")
	}
	mysqldbname, err := config.Config.GetString("default", "mysql_dbname")
	if err != nil {
		logging.Log("Problem with config file! (mysql_dbname)")
	}

	dbm, err := gorm.Open("mysql", mysqluser+":"+mysqlpass+
		"@("+mysqlhost+":"+mysqlport+")/"+mysqldbname+"?charset=utf8&parseTime=True")

	if err != nil {
		panic("Unable to connect to the database")
	} else {
		logging.Log("Database connection established.")
	}
	db = &dbm
	dbm.DB().Ping()
	dbm.DB().SetMaxIdleConns(10)
	dbm.DB().SetMaxOpenConns(100)
	db.LogMode(false)

	if !dbm.HasTable(&models.Node{}) {
		logging.Log("Node table not found, creating it now")
		dbm.CreateTable(&models.Node{})
	}

	if !dbm.HasTable(&models.Application{}) {
		logging.Log("Application table not found, creating it now")
		dbm.CreateTable(&models.Application{})
	}

	if !dbm.HasTable(&models.Container{}) {
		logging.Log("Container table not found, creating it now")
		dbm.CreateTable(&models.Container{})
	}

}

//Node database functions

//Create a new node in the database
func CreateNode(n *models.Node) (bool, error) {
	logging.Log("Creating Node: " + n.Hostname)

	//TODO: Check for auth

	err := db.Create(&n).Error
	if err != nil {
		return false, err
	}

	return true, err
}

//delete Node
func DeleteNode(id int64) (bool, error) {
	logging.Log("Deleting Node: ", id)

	node := models.Node{}

	err := db.Where(&models.Node{Id: id}).First(&node).Error

	if err != nil {
		return false, err
	}
	//  TODO: Check for auth
	//      Migrate all containers

	//Delete node from database
	err = db.Delete(&node).Error

	if err != nil {
		return false, err
	}

	return true, err

}

//Get node information
func GetNode(id int64) (*models.Node, error) {
	node := &models.Node{}

	err := db.Where(&models.Node{Id: id}).First(&node).Error

	return node, err
}

func GetNodes() ([]models.Node, error) {
	//Returns a list of all applications
	nodes := []models.Node{}

	err := db.Find(&nodes).Error

	return nodes, err
}

//Update node
func UpdateNode(node *models.Node) (bool, error) {

	newnode := models.Node{}
	err := db.Where(&models.Node{Id: node.Id}).First(&newnode).Error

	if err != nil {
		return false, err
	}
	err = db.Save(&node).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

//Application database functions

//Create a new application in the database
func CreateApplication(a *models.Application) (bool, error) {
	logging.Log("Creating Application: " + a.Name)

	//TODO: Check for auth

	err := db.Create(&a).Error
	if err != nil {
		return false, err
	}

	return true, err
}

//Get application information
func GetApplication(id int64) (*models.Application, error) {
	app := &models.Application{}

	err := db.Where(&models.Application{Id: id}).First(&app).Error

	return app, err
}

func GetApplications() ([]models.Application, error) {
	//Returns a list of all applications
	apps := []models.Application{}

	err := db.Find(&apps).Error

	return apps, err
}

//delete application from database
func DeleteApplication(id int64) (bool, error) {
	logging.Log("Deleting Application: ", id)

	app := models.Application{}

	err := db.Where(&models.Application{Id: id}).First(&app).Error

	if err != nil {
		return false, err
	}
	//  TODO: Check for auth
	//      Delete all containers

	//Delete application from database
	err = db.Delete(&app).Error

	if err != nil {
		return false, err
	}

	return true, err

}

//Update Application
func UpdateApplication(app *models.Application) (bool, error) {

	newapp := models.Application{}
	err := db.Where(&models.Application{Id: app.Id}).First(&newapp).Error

	if err != nil {
		return false, err
	}

	err = db.Save(&app).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

//Container database Functions

//Create a new node in the database
func CreateContainer(c *models.Container) (bool, error) {
	logging.Log("Creating Container: " + c.Name)

	//TODO: Check for auth

	err := db.Create(&c).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func UpdateContainer(cont *models.Container) (bool, error) {

	newcont := models.Container{}
	err := db.Where(&models.Application{Id: cont.Id}).First(&newcont).Error

	if err != nil {
		return false, err
	}

	err = db.Save(&cont).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

//Get container information
func GetContainer(id int64) (*models.Container, error) {
	cont := &models.Container{}

	err := db.Where(&models.Container{Id: id}).First(&cont).Error

	return cont, err
}

func GetContainers() ([]models.Container, error) {
	//Returns a list of all applications
	conts := []models.Container{}

	err := db.Find(&conts).Error

	return conts, err
}

func DeleteContainer(id int64) error {
	logging.Log("Deleting Application: ", id)

	cont := models.Container{}

	err := db.Where(&models.Container{Id: id}).First(&cont).Error

	if err != nil {
		return err
	}
	//  TODO: Check for auth
	//      Delete all containers

	//Delete application from database
	err = db.Delete(&cont).Error

	return err

}
