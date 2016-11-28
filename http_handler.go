package main

import "net/http"

func httpHandler(function http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		function(w, r)
	}
}
