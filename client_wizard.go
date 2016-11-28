package main

import "net/http"

func ClientWizardHandler(w http.ResponseWriter, r *http.Request) {
	// Render dashboard
	respond(w, r, http.StatusOK, "tmpl/content/newClientWizard.tmpl", nil, nil)
}