package database

import (
 	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/superordinate/kDaemon/models"
    "github.com/superordinate/kDaemon/logging"
    "github.com/superordinate/kDaemon/config"
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

    dbm, err := gorm.Open("mysql", mysqluser+ ":" + mysqlpass + 
    		"@(" + mysqlhost + ":" + mysqlport + ")/kdaemon?charset=utf8&parseTime=True")

    if(err != nil){
        panic("Unable to connect to the database")
    } else {
    	logging.Log("Database connection established.")
    }
    db = &dbm
    dbm.DB().Ping()
    dbm.DB().SetMaxIdleConns(10)
    dbm.DB().SetMaxOpenConns(100)
    db.LogMode(false)
 
    if !dbm.HasTable(&models.Node{}){
    	logging.Log("Node table not found, creating it now")
        dbm.CreateTable(&models.Node{})
    } 

    if !dbm.HasTable(&models.Application{}){
        logging.Log("Application table not found, creating it now")
        dbm.CreateTable(&models.Application{})
    } 

    if !dbm.HasTable(&models.Container{}){
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
     if  err != nil {
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
     if  err != nil {
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
     if  err != nil {
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

/* OLD CODE THAT MAY BE USEFUL

HELPER FUNCTIONS */


// //strips all whitespace out of a string
// func stripSpaces(str string) string {
//     return strings.Map(func(r rune) rune {
//         if unicode.IsSpace(r) {
//             // if the character is a space, drop it
//             return -1
//         }
//         // else keep it in the string
//         return r
//     }, str)
// }



// //Check if Username exists, returns false if username not taken
// func CheckForExistingUsername(u *models.User) bool {
// 	newUser := &models.User{}

// 	db.Where(&models.User{Username: u.Username}).First(&newUser)

// 	return newUser.Id == 0
// } 

// //Check if Email exists, returns false if email not taken
// func CheckForExistingEmail(u *models.User) bool {
// 	newUser := &models.User{}

// 	db.Where(&models.User{Email: u.Email}).First(&newUser)

// 	return newUser.Id == 0
// }

// //Check if passwords match
// func CheckForMatchingPassword(u *models.User) bool {
// 	newUser := &models.User{}

// 	db.Where("username = ?", u.Username).First(&newUser)
// 	fmt.Println(newUser)

// 	return newUser.Password == u.Password
// } 

// //Get All Users
// func GetUsers(u *[]models.User) {
// 	fmt.Println("Getting list of all applications")

// 	//Returns a list of all applications 
// 	userlist := []models.User{}

// 	db.Find(&userlist)
// 	//makes the list externally available
// 	*u = userlist
// }

// //Get User
// func GetUserByUsername(username string) *models.User {
// 	newUser := &models.User{}

// 	db.Where(&models.User{Username: username}).First(&newUser)

// 	return newUser
// }

// //UpdateUser
// func UpdateUser(u *models.User) {
// 	db.Save(&u)
// } 


// /* APPLICATION DATABASE THINGS */

// //Check if application exists in database, returns true if app exists
// func CheckApplicationExists(appname string) bool {
// 	if appname == "" {
// 		return true
// 	}

// 	newapp := &models.Application{}

// 	db.Where(&models.Application{Name: appname}).First(&newapp)

// 	return newapp.Id != 0
// }

// //Create application in DB
// func CreateApplication(a *models.Application) {
// 	fmt.Println("Creating Application: " + a.Name)
// 	db.Create(&a)

// }

// func UpdateApplication(a *models.Application) {
// 	fmt.Println("Updating Application: " + a.Name)
// 	db.Save(&a)
// }

// func DeleteApplication(a *models.Application, username string) {
// 	fmt.Println("Deleting Application: " + a.Name + " and all running instances")

// 	user := GetUserByUsername(username)
	
// 	if (user.Role == "admin") {
// 		//Get all running instances
// 		runningapps := []models.RunningApplication{}
		

// 		db.Where("application_id = ?", a.Id).Find(&runningapps)
// 		fmt.Println("running apps: ", runningapps)

// 		for _, element := range runningapps {
// 			DeleteRunningApplication(username, element.Name)
// 		}

// 		//Delete all dependencies
// 		dependency := models.Dependency{}
// 		db.Where("application_id = ?", a.Id).Delete(&dependency)

// 		//Delete all environment variables
// 		envars := []models.EnvironmentVariable{}
// 		db.Where("application_id = ?", a.Id).Delete(&envars)

// 		//Finally, delete the application
// 		db.Delete(&a)
	
// 	}

// }

// //Get Application List
// func GetApplications(a *[]models.Application) {
// 	fmt.Println("Getting list of all applications")

// 	//Returns a list of all applications 
// 	applicationList := []models.Application{}

// 	db.Find(&applicationList)
// 	//makes the list externally available
// 	*a = applicationList
// }

// //Get application by name
// func GetApplicationByName(appname string) *models.Application {

// 	newapp := &models.Application{}
// 	dependencies := []models.Dependency{}
// 	envvariables := []models.EnvironmentVariable{}

// 	//Get an application -- this doesnt grab associated dbs
// 	db.Where(&models.Application{Name: appname}).First(&newapp)

// 	db.Model(&newapp).Related(&dependencies)
// 	db.Model(&newapp).Related(&envvariables)

// 	fmt.Println(newapp)
// 	fmt.Println(dependencies)
// 	fmt.Println(envvariables)

// 	newapp.Dependencies = dependencies
// 	newapp.EnvironmentVariables = envvariables
// 	//Get dependencies

// 	return newapp
// }

// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// const (
//     letterIdxBits = 6                    // 6 bits to represent a letter index
//     letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
//     letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
// )

// func RandString(strlen int) string {
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
// 	result := make([]byte, strlen)
// 	for i := 0; i < strlen; i++ {
// 		result[i] = chars[rand.Intn(len(chars))]
// 	}
// 	return string(result)
// }

// /* USER APPLICATION THINGS */

// func AddRunningApplication(a *models.RunningApplication) {
// 	fmt.Println("Adding running Application: " + a.Name)

// 	db.Create(&a)
// }


// func GetRunningApplications(a *[]models.RunningApplication) {

// 	fmt.Println("Getting list of all running applications")

// 	//Returns a list of all applications 
// 	applicationList := []models.RunningApplication{}

// 	db.Find(&applicationList)
// 	LoadLogoForRunningApplications(&applicationList)
// 	//makes the list externally available
// 	*a = applicationList

// }

// func GetRunningApplicationsForUser(a *[]models.RunningApplication, u *models.User) {

// 	runningapps := []models.RunningApplication{}

// 	db.Where("owner = ?", u.Id).Find(&runningapps)

// 	LoadLogoForRunningApplications(&runningapps)
	
// 	*a = runningapps
// }

// func LoadLogoForRunningApplications(a *[]models.RunningApplication) {

// 	for index := range *a {
// 		application := models.Application{}

// 		db.Where("id = ?", (*a)[index].ApplicationID).First(&application)

// 		(*a)[index].Logo = application.Logo
// 	}
// }

// func GetRunningApplicationByName(name string) *models.RunningApplication{
// 	application := models.RunningApplication{}

// 	db.Where("name = ?", name).First(&application)

// 	return &application
// }

// func UpdateRunningApplication(a *models.RunningApplication) {
// 	db.Save(&a)
// }

// func DeleteRunningApplication(username, name string) bool{

// 	user := GetUserByUsername(username)
// 	application := GetRunningApplicationByName(name)


// 	if (user.Id == application.Owner || user.Role == "admin") {
// 		//mark as not running
// 		application.IsRunning = false
// 		UpdateRunningApplication(application)

// 		//remove from marathon
// 		url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps/" + name

// 		req, err := http.NewRequest("DELETE", url, nil)

// 		//Make the request
// 		res, err := http.DefaultClient.Do(req)

// 		if err != nil {
// 			panic(err) //Something is wrong while sending request
// 			}

// 		if res.StatusCode != 201 {
// 			fmt.Printf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
// 		}

// 		//Check if app is still running
// 		running := false

// 		for running {
// 			//if app is still running, don't do anything yet
// 			running, _ = CheckMarathonForRunningStatus(name)
// 		}

// 		//app isn't running, remove from database
// 		db.Delete(&application)

// 		return true
// 	} else {
// 		return false
// 	}
// }

// func CheckMarathonForRunningStatus(name string) (running bool, application models.MarathonApplication) {

// 	running =false
// 	url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps/" + name
// 	fmt.Println("Checking if application " + name + "is running yet.")
	
// 	time.Sleep(2 * time.Second)

// 	req, err := http.NewRequest("GET", url, nil)

// 	//Make the request
// 	res, err := http.DefaultClient.Do(req)

// 	if err != nil {
//     	panic(err) //Something is wrong while sending request
//  	}

// 	if res.StatusCode != 201 {
// 		fmt.Printf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
// 	}


// 	marathonapp := models.MarathonApplication{}


// 	body,err := ioutil.ReadAll(res.Body)
// 	if err != nil {
//     	panic(err) //Something is wrong while sending request
//  	}

// 	err = json.Unmarshal(body, &marathonapp)


// 	if err != nil {			
		
// 		panic(err)
// 		return
// 	}

// 	if (marathonapp.App.TasksRunning != 0) {
// 		running = true
// 	}
// 	application = marathonapp

// 	return 
// }
