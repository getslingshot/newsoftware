package main

import (
	"net/http"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdateClientHandler http handler to Update client
func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	ID := r.FormValue("id")

	var client Client
	//Grab values from form
	client.CompanyName = r.FormValue("company-name")
	client.Address = r.FormValue("company-address")
	client.Phone = r.FormValue("company-phone")
	client.ClientMessage = r.FormValue("company-message")

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("clients")

	//Build query to update company information
	query := bson.M{"_id": bson.ObjectIdHex(ID)}
	data := bson.M{
		"$set": bson.M{
			"companyName":   client.CompanyName,
			"address":       client.Address,
			"phone":         client.Phone,
			"clientMessage": client.ClientMessage,
		},
	}

	err := c.Update(query, data)

	if err != nil {
		logger.Println("error on get client", err)
		return
	}

	logger.Println("Client Information updated successful")

	http.Redirect(w, r, "/clients", 302)
}
