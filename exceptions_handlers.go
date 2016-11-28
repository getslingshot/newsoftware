package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createExceptionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var exception Exception

	if err := decodeBody(r, &exception); err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	invalid, errs := exception.Validate()

	if invalid {
		e := sendError()
		e.Error = "Invalid payload"
		e.Errors = errs

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("exceptions")
	exception.ID = bson.NewObjectId()

	if err := c.Insert(exception); err != nil {
		e := sendError()
		e.Error = "Cant process request"

		logger.Println("Error on insert", err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	respondJSON(w, r, http.StatusOK, exception)
}

func deleteExceptionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("exceptions")
	query := bson.M{"_id": bson.ObjectIdHex(id)}

	if err := c.Remove(query); err != nil {
		e := sendError()
		e.Error = "Cant process request"

		logger.Println(err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	// return id, so Front End can delete html content
	d := map[string]string{"_id": id}
	respondJSON(w, r, http.StatusOK, d)
}

func updateExceptionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}
	var exception Exception

	if err := decodeBody(r, &exception); err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	invalid, errs := exception.Validate()

	if invalid {
		e := sendError()
		e.Error = "Invalid payload"
		e.Errors = errs

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("exceptions")
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	exception.ID = bson.ObjectIdHex(id)

	if err := c.Update(query, exception); err != nil {
		e := sendError()
		e.Error = "Cant process request"

		logger.Println(err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	respondJSON(w, r, http.StatusOK, exception.Public())
}
