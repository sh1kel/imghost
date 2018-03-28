package main

import (
	"net/http"
	"html/template"
	"fmt"
	"mime/multipart"
	"io/ioutil"
	"github.com/labstack/gommon/log"
	"encoding/json"
	"github.com/rainycape/unidecode"
	"path"
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
var DownloadHandler = http.HandlerFunc(DownloadData)


func RootRoute(w http.ResponseWriter, r *http.Request) {
	var user = User{Name: "Anonym"}
	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "http://localhost:8081/auth", 307)
	}
	if session.Values["username"] != nil {
		user.Name = session.Values["username"].(string)
	}
	indexTemplate, err := template.ParseFiles("templates/index.html")

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

func DownloadData(w http.ResponseWriter, r *http.Request) {
	var fName string
	type jsonResponse struct {
		FileUrl	string	`json:"FileUrl"`
		Status	string	`json:"Status"`
	}
	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["authenticated"] != true {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("<html>Unauthorized</html>"))
		return
	}
	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "{\"Error\":\"%v\"}", err)
		return
	}
	defer file.Close()
	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		fName, err = saveFile(w, file, handle)
	case "image/png":
		fName, err = saveFile(w, file, handle)
	default:
		responseData := jsonResponse{FileUrl: "File type unsuported"}
		data, err := json.Marshal(responseData)
		if err != nil {
			fmt.Fprintf(w, "{\"Status\":\"%v\"}", err)
			return
		}
		w.WriteHeader(406)
		w.Write(data)
	}
	responseData := jsonResponse{FileUrl: baseUrl + fName, Status: "ok"}
	data, err := json.Marshal(responseData)
	if err != nil {
		fmt.Fprintf(w, "{\"Status\":\"%v\"}", err)
		return
	}
	w.Write(data)
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) (string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return "", err
	}
	filename := path.Base(handle.Filename)
	filename = unidecode.Unidecode(filename)
	err = ioutil.WriteFile("./upload/"+filename, data, 0664)
	log.Printf("Saving %s\n", filename)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return "", err
	}
	return filename, nil
}