package models

import (
	"regexp"
	"unicode"
	
)

//User struct for database entries
type User struct {
  Id       	int64
  Username 	string	`sql:"size:255; not null; unique;"`
  Email		string 	`sql:"size:255; not null; unique;"` 
  FirstName	string	`sql:"size:30; not null;"`
  Surname 	string 	`sql:"size:30; not null;"`   
  Password 	string 	`sql:"size:255; not null;"`
  Role 		string 	`sql:"size:30"` 
  IsEnabled	bool 	`sql:"default:true"`
  	ConfirmPassword 	string 	`sql:"-"`	//These dont get put in the database
  	Message				string 	`sql:"-"`	//These dont get put in the database
}

//Validates inputs for registration page
func (u *User) ValidateRegister() {


		errorstring := "";


		if VerifyName(u.Username) {
			errorstring = errorstring + "Username is invalid. A-Za-z0-9 only. -- "
		}
		if VerifyName(u.FirstName) {
			errorstring = errorstring + "First Name is invalid. A-Za-z0-9 only. -- "
		}
		if VerifyName(u.Surname) {
			errorstring = errorstring + "Last Name is invalid. A-Za-z0-9 only. -- "
		}
	


		if !VerifyEmail(u.Email) {
			errorstring = errorstring + "Email is invalid. --"
		}

		tenOrMore, _, _, _ := VerifyPassword(u.Password)
		if (!(tenOrMore)) {
			errorstring = errorstring + "Password must be at least 10 characters long"
		}

		if (u.Password != u.ConfirmPassword) {
			errorstring = errorstring + "Passwords don't match. -- "
		}

		u.Message = errorstring

}

//Validates inputs for login page
func (u *User) ValidateLogin() {
	errorstring := "";

	if VerifyName(u.Username) {
		errorstring = errorstring + "Username is invalid. A-Za-z0-9 only. -- "
	}

	tenOrMore, _, _, _ := VerifyPassword(u.Password)
	if (!(tenOrMore)) {
		errorstring = errorstring + "Password must be at least 10 characters long"
	}

	u.Message = errorstring


}

//Validates inputs for updating password on profile page
func (u *User) ValidateNewPassword() {
	errorstring := "";

	tenOrMore, _, _, _ := VerifyPassword(u.Password)
	if (!(tenOrMore)) {
		errorstring = errorstring + "Password must be at least 10 characters long"
	}

	if (u.Password != u.ConfirmPassword) {
		errorstring = errorstring + "Passwords don't match. -- "
	}

	u.Message = errorstring


}

//returns true if invalid
func VerifyName(s string) bool {
	reg, err := regexp.Compile(`\W`)

	if err != nil {
		panic(err)
	}

	return reg.MatchString(string(s))

}

//return true if valid
func VerifyEmail(s string) bool {

	reg , err := regexp.Compile(`\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3}`)

	if err != nil {
		panic(err)
	}

	return reg.MatchString(string(s))
}

//returns true when it contains each requirement.
func VerifyPassword(s string) (tenOrMore, number, upper, special bool) {
    letters := 0
    for _, s := range s {
        switch {
        case unicode.IsUpper(s):
            upper = true
            letters++
        case unicode.IsPunct(s) || unicode.IsSymbol(s):
            special = true
        case unicode.IsLetter(s) || unicode.IsNumber(s):
            letters++
        default:
            //return false, false, false, false
        }
    }
    tenOrMore = letters >= 10
    return
}