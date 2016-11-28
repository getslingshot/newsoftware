package main

import (
	"net/http"
	"os"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RequestPasswordHandler http handler
func RequestPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	e := r.FormValue("email")

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")
	query := bson.M{"email": e}

	var user User
	err := c.Find(query).One(&user)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	user.RecoveryToken = bson.NewObjectId().Hex()

	err = c.Update(query, user)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	url := os.Getenv("RECOVERY_PASSWORD")

	link := `To change password follow this <a href=` + url + user.RecoveryToken + `>link</a>`

	to := make([]string, 1)
	to[0] = user.Email

	data := Email{
		Recipients: to,
		Subject:    "Slingshot recovery password instructions",
		Message:    link,
	}

	err = email(data)

	if err != nil {
		logger.Println("Error on send email", err)
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
