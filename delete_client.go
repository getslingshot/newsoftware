package main

import (
	"net/http"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdateClientHandler http handler to Update client
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	ID := r.FormValue("id")

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("clients")

	//Build query to update company information
	query := bson.M{"_id": bson.ObjectIdHex(ID)}

	err := c.Remove(query)

	if err != nil {
		logger.Println("error on DELETE client", err)
		return
	}

	logger.Println("Client Information deleted successful")

	http.Redirect(w, r, "/clients", 302)
}
