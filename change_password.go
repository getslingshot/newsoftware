package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ChangePasswordHandler http handler
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	password := r.FormValue("password")
	confirmPass := r.FormValue("confirm-password")
	token := r.FormValue("token")

	if password != confirmPass {
		http.Redirect(w, r, "/cp?t="+token, http.StatusPermanentRedirect)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")

	query := bson.M{
		"recoveryToken": token,
	}

	var user User
	err := c.Find(query).One(&user)

	if err != nil {
		logger.Println("error on query for user with token", err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	user.RecoveryToken = ""
	p, err := Crypt([]byte(password))

	if err != nil {
		// TODO: Display error
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}
	user.Password = p

	err = c.Update(query, user)

	if err != nil {
		// TODO: Redirect or display error
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

// ChangePassFormHandler http handler
func ChangePassFormHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")
	token := r.URL.Query().Get("t")

	query := bson.M{
		"recoveryToken": token,
	}

	var user User
	err := c.Find(query).One(&user)

	if err != nil {
		// TODO
		// Send proper template error here
		logger.Println("Background Search Template Parse Error: ", err)
		return
	}

	tmpl, err := template.ParseFiles(
		"tmpl/content/change_password.html",
	)

	if err != nil {
		// TODO
		// Send proper template error here
		logger.Println("Background Search Template Parse Error: ", err)
		return
	}

	err = tmpl.Execute(w, user)
}
