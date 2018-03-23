package main

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"fmt"
)

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

var SignInHandler = http.HandlerFunc(SignIn)
var SignUpHandler = http.HandlerFunc(SignUp)
var LogoutHandler = http.HandlerFunc(Logout)

func SignIn(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := db.QueryRow("select password from auth where username=$1", creds.Username)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Write([]byte("Auth failed"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		w.Write([]byte("Auth failed"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = 3600
	session.Values["authenticated"] = true
	session.Values["username"] = creds.Username

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if session.Values["authenticated"] != true {
		w.Write([]byte("Forbidden"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("insert into auth (username, password) values ($1, $2)", creds.Username, string(hashedPassword))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := fmt.Sprintf("User created %s: %s", creds.Username, string(hashedPassword))
	w.Write([]byte(user))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "imghost-cookie")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if session.Values["authenticated"] != true {
		w.Write([]byte("Forbidden"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//w.Write([]byte("Logged out"))
	http.Redirect(w, r, "http://localhost:8081", 307)

}
