package main

import "net/http"

func LogoffHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}
