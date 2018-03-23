package main

import (
	"net/http"
	"html/template"
)

type User struct {
	Name	string

}

var DumbHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dumb"))
})

var RootHandler = http.HandlerFunc(RootRoute)
var AuthHandler = http.HandlerFunc(AuthRoute)
var NotFoundHandler = http.HandlerFunc(NotFoundRoute)


func RootRoute(w http.ResponseWriter, r *http.Request) {
	var user = User{Name: "Preved"}
	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "http://localhost:8081/auth", 307)
	}
	//w.Write([]byte("This is /"))
	user.Name = session.Values["username"].(string)
	//tmpl := template.New("index")
	//template.Must(tmpl.ParseFiles("templates/index.html"))
	indexTemplate, err := template.ParseFiles("templates/index.html")
	//template.Must(tmpl.Parse("Hello, {{.Name}}"))

	err = indexTemplate.Execute(w, user)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


}

func AuthRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/auth.html"))
	tmpl.Execute(w, nil)
}

func NotFoundRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Nothing found..."))
}