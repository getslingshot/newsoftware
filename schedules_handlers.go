package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var collection Schedules

	if err := decodeBody(r, &collection); err != nil {
		respondJSON(w, r, http.StatusBadRequest, err)
		return
	}

	e := sendError()

	for i, doc := range collection {
		invalid, errs := doc.Validate()

		if invalid {
			e.Error = "Invalid payload"

			for k, v := range errs {
				e.Errors[k+"-"+strconv.Itoa(i)] = v
			}
		}
	}

	if e.Error != "" {
		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	// let's sort collection, by default sort by schedule.Opens property
	sort.Sort(collection)

	// Let's validate ranges
	valid, errs := validateScheduleRange(collection)
	if !valid {
		e := sendError()
		e.Error = "Invalid payload"
		e.Errors = errs

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("schedules")

	// First let's remove all schedules.
	if _, err := c.RemoveAll(bson.M{}); err != nil && err != mgo.ErrNotFound {
		e := sendError()
		e.Error = "Cant process request right now"

		logger.Println(err)
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	// Then let's add new schedules to the collection
	bulk := c.Bulk()

	for _, s := range collection {
		s.ID = bson.NewObjectId()
		bulk.Insert(s)
	}

	if _, err := bulk.Run(); err != nil {
		e := sendError()
		e.Error = "Cant process request"
		logger.Println("Cant insert schedules", err)

		respondJSON(w, r, http.StatusUnprocessableEntity, e)
		return
	}

	data := make([]map[string]interface{}, len(collection)+6)

	for i, s := range collection {
		data[i] = s.Public()
	}

	respondJSON(w, r, http.StatusOK, data)
}
