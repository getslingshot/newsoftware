package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	// http://goo.gl/PIWVzT
	id := r.URL.Query().Get("id")

	if id == "" {
		logger.Println("User ID not provided")
		http.Redirect(w, r, "/", 302)
		return
	}

	// TODO:
	// Validate against current user account
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")

	var user User
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)

	if err != nil {
		logger.Println("error on get user", err)
		return
	}

	roles := []Role{}
	c = db.C("roles")
	err = c.Find(bson.M{}).All(&roles)

	for _, r := range roles {
		if r.ID == user.RoleID {
			user.RoleTitle = r.Title
			break
		}
	}

	if err != nil {
		logger.Println("error on get user", err)
		return
	}
	d := M{
		"user":  user,
		"roles": roles,
	}

	respond(w, r, http.StatusOK, "tmpl/content/edit_user.tmpl", d, nil)
}

// M http://stackoverflow.com/questions/25329647/golang-template-with-multiple-structs
type M map[string]interface{}
