package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ConfirmHandler allow user to enter password for new account based on
// VerifyToken properpty
func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")
	token := r.URL.Query().Get("t")

	query := bson.M{
		"verifyToken": token,
	}

	var user User
	err := c.Find(query).One(&user)

	if err != nil {
		logger.Println("error on query for user with token", err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	tmpl, err := template.ParseFiles(
		"tmpl/content/password.html",
	)

	if err != nil {
		logger.Println("Background Search Template Parse Error: ", err)
		// TODO
		// Send proper template error here
		return
	}

	err = tmpl.Execute(w, user)
}
