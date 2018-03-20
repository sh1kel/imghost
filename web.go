package main

import "net/http"

var DumbHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dumb"))
})

var RootHandler = http.HandlerFunc(RootRoute)

func RootRoute(w http.ResponseWriter, r *http.Request) {


}