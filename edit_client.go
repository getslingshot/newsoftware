package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func EditClientHandler(w http.ResponseWriter, r *http.Request) {
	// http://goo.gl/PIWVzT
	id := r.URL.Query().Get("id")

	if id == "" {
		logger.Println("Client ID not provided")
		http.Redirect(w, r, "/", 302)
		return
	}

	// TODO:
	// Validate against current user account
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("clients")

	var client Client
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&client)

	if err != nil {
		logger.Println("error on get client", err)
		return
	}

	d := M{
		"client": client,
		// 	"example": example,		//if you need multi-structs add them here
	}

	respond(w, r, http.StatusOK, "tmpl/content/edit_client.tmpl", d, nil)

}
