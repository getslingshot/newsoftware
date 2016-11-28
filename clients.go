package main

import (
	"net/http"
	"os"

	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Client struct associated to `clients` collection
type Client struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	CompanyName string        `bson:"companyName" json:"companyName"`
	Address     string        `bson:"address" json:"address"`
	City        string        `bson:"city" json:"city"`
	State       string        `bson:"state" json:"state"`
	Zipcode     string        `bson:"zipcode" json:"zipcode"`

	//Contact information
	LastName  string `bson:"lastName" json:"lastName"`
	FirstName string `bson:"firstName" json:"firstName"`
	Phone     string `bson:"phone" json:"phone"`
	Email     string `bson:"email" json:"email"`

	// Non input fields
	Slug          string `bson:"slug" json:"slug"`
	ClientMessage string `bson:"clientMessage" json:"clientMessage"`
}

// ClientsHandler http handler for GET /clients
func ClientsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getClients(w, r)
		return

	case http.MethodPost:
		createClients(w, r)
		return

	default:
		notSupported(w, r)
		return
	}
}

func getClients(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("clients")

	clients := []Client{}
	err := c.Find(bson.M{}).All(&clients)

	if err != nil {
		logger.Println("error on get clients", err)
		return
	}

	cu := context.Get(r, "currentUser").(User)

	d := M{
		"user":    cu,
		"clients": clients,
	}

	respond(w, r, http.StatusOK, "tmpl/content/clients.tmpl", d, nil)
}

// Validate validator for Client object
func (c *Client) Validate() (bool, map[string]string) {
	v := validator.NewValidator()
	v.Present("companyName", c.CompanyName)
	v.Present("firstName", c.FirstName)
	v.Present("lastName", c.LastName)
	v.Present("email", c.Email)
	v.Email("email", c.Email)

	return v.IsValid(), v.Errors
}

func createClients(w http.ResponseWriter, r *http.Request) {
	var client Client

	if err := decodeBody(r, &client); err != nil {
		respondJSON(w, r, 422, err)
		return
	}

	invalid, errors := client.Validate()

	if invalid {
		respondJSON(w, r, 422, errors)
		return
	}

	client.ID = bson.NewObjectId()
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("clients")
	err := c.Insert(client)

	if err != nil {
		logger.Println("error on inserting new client to database", err)
		return
	}

	var user User
	user.ID = bson.NewObjectId()
	user.ClientID = client.ID
	user.FirstName = client.FirstName
	user.LastName = client.LastName
	user.Phone = client.Phone
	user.Email = client.Email
	user.Verified = false
	user.VerifyToken = bson.NewObjectId().Hex()

	c = db.C("users")
	err = c.Insert(user)

	if err != nil {
		logger.Println("error on inserting new client to database", err)
		return
	}

	url := os.Getenv("CONFIRM_URL")

	link := `<a href=` + url + user.VerifyToken + `>validate url</a>`
	to := make([]string, 1)
	to[0] = user.Email

	data := Email{
		Recipients: to,
		Subject:    "Verify account",
		Message:    link,
	}

	err = email(data)

	if err != nil {
		logger.Println("Error on send email", err)
	}

	respondJSON(w, r, http.StatusCreated, client)
	return
}
