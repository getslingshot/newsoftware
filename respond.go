package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Respond struct to respond with data and errors
type Respond struct {
	Data   interface{}
	Errors interface{}
}

var funcs = template.FuncMap{
	"sum": sum,
}

func respond(w http.ResponseWriter, r *http.Request, status int, layout string, d interface{}, errors interface{}) {
	tmpl := template.Must(template.New("base.tmpl").Funcs(funcs).ParseFiles(
		"tmpl/layout2/base.tmpl",
		"tmpl/layout2/footer.tmpl",
		"tmpl/layout2/header.tmpl",
		layout,
	))

	var data Respond
	data.Data = d
	if errors == nil {
		data.Errors = make(map[string]interface{})
	} else {
		data.Errors = errors
	}

	err := tmpl.Execute(w, data)

	if err != nil {
		w.WriteHeader(status)
		// Send proper template error here
		fmt.Println("Background Search Template Execution Error: ", err)
		return
	}
}
