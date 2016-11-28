package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func currentUser() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var u User
			cookie, err := r.Cookie("session")
			cookieValue := make(map[string]string)

			email := ""

			if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
				email = cookieValue["email"]
			}

			u.Email = email

			db := context.Get(r, "db").(*mgo.Database)
			c := db.C("users")

			// Load Current User Into Memory
			err = c.Find(bson.M{"email": email}).One(&u)
			if err != nil {
				logger.Println("Could Not Locate Customer Record: ", err)
				header := r.Header.Get("Content-Type")
				// Respond for json requests
				if strings.Contains(header, "application/json") {
					respondJSON(w, r, http.StatusForbidden, nil)
					return
				}

				http.Redirect(w, r, "/", 302)
				return
			}
			_, err = r.Cookie("session")

			if err != nil {
				logger.Println("You Must Be Logged In")
				http.Redirect(w, r, "/", 302)
				return
			}

			context.Set(r, "currentUser", u)

			h.ServeHTTP(w, r)
		})
	}
}
