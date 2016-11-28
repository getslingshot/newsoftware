package main

import (
	"net/http"
	"strings"
)

func authenticate() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie("session")

			if err != nil {
				header := r.Header.Get("Content-Type")
				// Respond for json requests
				if strings.Contains(header, "application/json") {
					w.Header().Set("Content-Type", "application/json")
					e := sendError()
					e.Error = "You must have a session"
					respondJSON(w, r, http.StatusForbidden, e)
					return
				}

				logger.Println("You Must Be Logged In")
				http.Redirect(w, r, "/", 302)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
