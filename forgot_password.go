package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// ForgotPasswordHandler http handler
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"tmpl/content/forgot_password.html",
	)

	if err != nil {
		fmt.Println("Background Search Template Parse Error: ", err)
		// TODO
		// Send proper template error here
		return
	}

	err = tmpl.Execute(w, nil)
}
