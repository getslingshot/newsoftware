package main

import (
	"net/http"
	"os"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
)

func mongo(instance *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := instance.Copy()
			defer c.Close()

			context.Set(r, "db", c.DB(os.Getenv("MONGO_NAME")))

			h.ServeHTTP(w, r)
		})
	}
}
