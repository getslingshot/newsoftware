package main

import (
	"net/http"

	"github.com/gorilla/context"
	//	mgo "gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	//  //  //  *****Need to add a lookup user function!!!

	// users, err := users(r)

	// if err != nil {
	// 	logger.Println("error on get users", err)
	// 	return
	// }

	cu := context.Get(r, "currentUser").(User)

	d := M{
		"user": cu,
		//"users": users,
	}

	respond(w, r, http.StatusOK, "tmpl/content/user_profile.tmpl", d, nil)
}
