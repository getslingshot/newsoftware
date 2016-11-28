package main

import (
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoginHandler http handler
// Validate user credentials
// Verify computers: if any, redirect to getStarted, otherwise to /dashboard
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("login-email")
	password := r.FormValue("login-password")
	passwordByte := []byte(password)

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")

	user := User{}
	err := c.Find(bson.M{"email": userEmail}).One(&user)

	if err != nil {
		logger.Println("Mongo Email Find Error: ", err)
	}

	err = Compare(user.Password, passwordByte)

	if err != nil {
		logger.Println("Mongo Login Auth Error: ", err)
		http.Redirect(w, r, "/", 302)
		return
	}

	value := map[string]string{
		"email":     userEmail,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}

		http.SetCookie(w, &cookie)
	}

	http.Redirect(w, r, "/clients", 302)
}
