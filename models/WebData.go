/*
	WebData struct will be passed to webpages in the UI. It will handle things like User Status, User Data, and Page Data


*/

package models

import ()

type WebData struct {
	LoggedIn    bool        //This will be either logged in, or not.
	CurrentUser User        //User data for logged in user
	Message     string      //How we will pass messages to the user
	PageData    interface{} //Data that we want to pass to page.
}
