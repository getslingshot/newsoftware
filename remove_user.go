package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RemoveUserHandler http handler to remove / delete user
func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {

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

	err := c.Remove(query)

	if err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	d := make(map[string]string, 1)
	d["_id"] = id

	respondJSON(w, r, http.StatusOK, d)
}
