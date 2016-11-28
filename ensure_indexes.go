package main

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

func ensureIndexes(session *mgo.Session) {

	index := mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}

	d := session.Copy()
	defer d.Close()

	c := d.DB(os.Getenv("MONGO_NAME")).C("users")
	err := c.EnsureIndex(index)

	if err != nil {
		logger.Println("can't ensure index to users collection", err)
		// TODO: email panic alert
		return
	}
}
