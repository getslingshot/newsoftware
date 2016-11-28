package main

import (
	"net/http"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdateUserHandler http handler to Update user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	// http://goo.gl/PIWVzT
	id := r.URL.Query().Get("id")

	if id == "" {
		e := sendError()
		e.Error = "Valid ID required"

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")
	query := bson.M{"_id": bson.ObjectIdHex(id)}

	var user User
	if err := decodeBody(r, &user); err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	var found User
	err := c.Find(query).One(&found)

	if err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	found.Suspended = user.Suspended

	if user.RoleID != "" {
		found.RoleID = user.RoleID
	}

	err = c.Update(query, found)

	if err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	d := map[string]string{"_id": id}

	respondJSON(w, r, http.StatusOK, d)
}
