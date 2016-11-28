package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createConditionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var condition Condition

	if err := decodeBody(r, &condition); err != nil {
		e := sendError()
		e.Error = "Invalid payload"

		respondJSON(w, r, http.StatusBadRequest, e)
		return
	}

	err := condition.Validate()

	if err.Error != "" {
		e := sendError()
		e.Errors = err.Errors
		e.Error = "Invalid payload"

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("conditions")
	condition.ID = bson.NewObjectId()

	if err := c.Insert(condition); err != nil {
		e := sendError()
		e.Error = "Cant process request"
		logger.Println("Cant insert condition", err)

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	respondJSON(w, r, http.StatusOK, condition)
}

func getConditionsIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var collection Conditions
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("conditions")

	if err := c.Find(bson.M{}).All(&collection); err != nil {
		e := sendError()
		e.Error = "Cant process request"
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	d := make([]map[string]interface{}, len(collection))

	for i, c := range collection {
		d[i] = c.Public()
	}

	respondJSON(w, r, http.StatusOK, d)
}

func deleteConditionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("conditions")
	query := bson.M{"_id": bson.ObjectIdHex(id)}

	if err := c.Remove(query); err != nil {
		if err == mgo.ErrNotFound {
			e := sendError()
			e.Error = "Resource not found"
			respondJSON(w, r, http.StatusNotFound, e)
			return
		}

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

func updateConditionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}
	var condition Condition

	if err := decodeBody(r, &condition); err != nil {
		e := sendError()
		e.Error = "Invalid payload"

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	err := condition.Validate()

	if err.Error != "" {
		e := sendError()
		e.Error = "Invalid payload"
		e.Errors = err.Errors

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("conditions")
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	condition.ID = bson.ObjectIdHex(id)

	if err := c.Update(query, condition); err != nil {
		e := sendError()
		e.Error = "Cant process request"

		logger.Println(err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	respondJSON(w, r, http.StatusOK, condition.Public())
}
