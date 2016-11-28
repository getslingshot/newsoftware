package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// RootHandler http handler to display login page
func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"index.html",
	)

	if err != nil {
		fmt.Println("Background Search Template Parse Error: ", err)
		// TODO
		// Send proper template error here
		return
	}

	err = tmpl.Execute(w, nil)
}
