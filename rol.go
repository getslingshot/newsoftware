package main

import "gopkg.in/mgo.v2/bson"

// Role model
type Role struct {
	ID    bson.ObjectId `bson:"_id" json:"_id"`
	Title string        `bson:"title" json:"title"`
}
