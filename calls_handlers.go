package main

import (
	"net/http"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createCallHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var call Call

	if err := decodeBody(r, &call); err != nil {
		respondJSON(w, r, http.StatusBadRequest, err)
		return
	}

	invalid, errs := call.Validate()
	if invalid {
		e := sendError()
		e.Errors = errs
		e.Error = "Invalid payload"

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("calls")
	call.ID = bson.NewObjectId()

	if err := c.Insert(call); err != nil {
		e := sendError()
		e.Error = "Cant process request"
		logger.Println("Cant insert call", err)

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	respondJSON(w, r, http.StatusOK, call)
}

func getCallsIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var collection Calls
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("calls")

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

func deleteCallHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("calls")
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

func updateCallHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" || !bson.IsObjectIdHex(id) {
		respondJSON(w, r, http.StatusNotFound, nil)
		return
	}
	var call Call

	if err := decodeBody(r, &call); err != nil {
		respondJSON(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	invalid, errs := call.Validate()

	if invalid {
		e := sendError()
		e.Error = "Invalid payload"
		e.Errors = errs

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("calls")
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	call.ID = bson.ObjectIdHex(id)

	if err := c.Update(query, call); err != nil {
		e := sendError()
		e.Error = "Cant process request"

		logger.Println(err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	respondJSON(w, r, http.StatusOK, call.Public())
}
