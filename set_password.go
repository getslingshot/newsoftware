package main

import (
	"net/http"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PasswordHandler allow user to enter password for new account based on
// VerifyToken properpty
func PasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	password := r.FormValue("password")
	confirmPass := r.FormValue("confirm-password")
	token := r.FormValue("token")

	if password != confirmPass {
		http.Redirect(w, r, "/set-pass", http.StatusPermanentRedirect)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")

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

	user.VerifyToken = ""
	user.Verified = true
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
