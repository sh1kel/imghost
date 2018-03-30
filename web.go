package main

import (
	"net/http"
	"html/template"
	"fmt"
	"encoding/json"
	"log"
	"crypto/md5"
	"encoding/hex"
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
var UploadHandler = http.HandlerFunc(UploadData)
var FilesHandler = http.HandlerFunc(FilesRoute)

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

func UploadData(w http.ResponseWriter, r *http.Request) {
	var fName string
	var fSize int64
	type jsonResponse struct {
		Name	string	`json:"name"`
		Size	int64	`json:"size"`
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
	userName := session.Values["username"].(string)
	userNameHash := md5.New()
	userNameHash.Write([]byte(userName))
	hashedUserName := hex.EncodeToString(userNameHash.Sum(nil))
	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "{\"Error\":\"%v\"}", err)
		return
	}
	defer file.Close()
	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		fName, fSize, err = saveFile(file, handle, "jpg", hashedUserName)
	case "image/png":
		fName, fSize, err = saveFile(file, handle, "png", hashedUserName)
	case "image/gif":
		fName, fSize, err = saveFile(file, handle, "gif", hashedUserName)
	default:
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(406)
		return
	}
	responseData := jsonResponse{Name: baseUrl + hashedUserName + "/" + fName, Size: fSize}
	data, err := json.Marshal(responseData)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(data)
}

func FilesRoute(w http.ResponseWriter, r *http.Request) {
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
	userName := session.Values["username"].(string)
	userNameHash := md5.New()
	userNameHash.Write([]byte(userName))
	userName = hex.EncodeToString(userNameHash.Sum(nil))
	files := scanUploads(uploadDir, userName)
	response, err := json.Marshal(files.items)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(response)
}

