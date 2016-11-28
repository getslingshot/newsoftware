package main

import (
	"net/http"
	"os"

	"github.com/gorilla/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User model associated to `users` collection
type User struct {
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	ClientID      bson.ObjectId `bson:"clientId,omitempty" json:"clientId,omitempty"`
	RoleID        bson.ObjectId `bson:"roleId,omitempty" json:"roleId,omitempty"`
	FirstName     string        `bson:"name" json:"name"`
	LastName      string        `bson:"lastName" json:"lastName"`
	Email         string        `bson:"email" json:"email"`
	Phone         string        `bson:"phone" json:"phone"`
	Location      string        `bson:"location" json:"location"`
	Position      string        `bson:"position" json:"position"`
	Notes         string        `bson:"notes" json:"notes"`
	Password      []byte        `bson:"password" json:"password"`
	Verified      bool          `bson:"verified" json:"verified"`
	Suspended     bool          `bson:"suspended" json:"suspended"`
	RecoveryToken string        `bson:"recoveryToken" json:"recoveryToken"`
	VerifyToken   string        `bson:"verifyToken" json:"verifyToken"`
	RoleTitle     string
}

// UsersHandler http handler
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
		return

	case http.MethodGet:
		getUsers(w, r)
		return

	case http.MethodDelete:
		RemoveUserHandler(w, r)
		return

	case http.MethodPut:
		UpdateUserHandler(w, r)
		return

	default:
		notSupported(w, r)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Println("error on parse form", err)
		return
	}

	var user User
	user.ID = bson.NewObjectId()
	user.FirstName = r.FormValue("name")
	user.Notes = r.FormValue("notes")
	user.Position = r.FormValue("position")
	user.Email = r.FormValue("email")
	user.Phone = r.FormValue("phone")
	user.Location = r.FormValue("location")
	user.LastName = r.FormValue("lastName")
	user.VerifyToken = randomToken(8)
	user.Verified = false

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")
	err := c.Insert(user)

	if err != nil {
		if mgo.IsDup(err) {
			u, _ := users(r)
			e := make(map[string]interface{})
			e["email"] = "Email already exist"

			respond(w, r, 422, "tmpl/content/users.tmpl", u, e)
			return
		}

		logger.Println("error on inserting new user to database", err)
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

	// Since we are using same url /users for GET / POST we just run GET
	// handler
	getUsers(w, r)
	return
}

//func user()

func users(r *http.Request) ([]User, error) {
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("users")

	users := []User{}
	err := c.Find(bson.M{}).All(&users)

	return users, err
}


func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := users(r)

	if err != nil {
		logger.Println("error on get users", err)
		return
	}

	cu := context.Get(r, "currentUser").(User)

	d := M{
		"user": cu,
		"users": users,
	}

	respond(w, r, http.StatusOK, "tmpl/content/users.tmpl", d, nil)
}