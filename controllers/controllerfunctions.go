package controllers

import (
 	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/superordinate/kDaemon/models"
    "fmt"
    "os"

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

	fmt.Println("Initializing Database connection.")
	mysqlhost := os.Getenv("MYSQL_HOST")
	mysqluser := os.Getenv("MYSQL_USER")
	mysqlpass := os.Getenv("MYSQL_PASSWORD")

	fmt.Println("mysql", mysqluser+ ":" + mysqlpass + 
    		"@(" + mysqlhost + ")/kdaemon")
    dbm, err := gorm.Open("mysql", mysqluser+ ":" + mysqlpass + 
    		"@(" + mysqlhost + ")/kdaemon?charset=utf8&parseTime=True")

	fmt.Println("Doing a thing.")
    if(err != nil){
        panic("Unable to connect to the database")
    } else {
    	fmt.Println("Database connection established.")
    }
    fmt.Println("Doing a thing.")
    db = &dbm
    dbm.DB().Ping()
    dbm.DB().SetMaxIdleConns(10)
    dbm.DB().SetMaxOpenConns(100)
    db.LogMode(true)
 
    if !dbm.HasTable(&models.Node{}){
    	fmt.Println("Node table not found, creating it now")
        dbm.CreateTable(&models.Node{})
    } 
    fmt.Println("Doing a thing.")
}

//Node database functions
//Create a new node in the database
func CreateNode(n *models.Node) (bool, error) {
     fmt.Println("Creating Node: " + n.Hostname)

     //TODO: Check for auth

     err := db.Create(&n).Error
     if  err != nil {
        return false, err
     }

     return true, err
}

func DeleteNode(id int64) (bool, error) {
    fmt.Println("Deleting Node: ", id)

    node := models.Node{}

    err := db.Where(&models.Node{Id: id}).First(&node).Error 

    if err != nil {
        return false, err
    }
    //TODO: Check for auth
    //      Migrate all containers

    //Delete node from database
    err = db.Delete(&node).Error

    if err != nil {
       return false, err
    }

    return true, err

}


/* HELPER FUNCTIONS */


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
