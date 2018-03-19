package main

import (
	"encoding/json"
	"database/sql"
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"fmt"
	"github.com/labstack/gommon/log"
)

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

var SigninHandler = http.HandlerFunc(Signin)
var SignupHandler = http.HandlerFunc(Signup)


func Signin(w http.ResponseWriter, r *http.Request){
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
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(`{"Auth": "ok"}`))

}

func Signup(w http.ResponseWriter, r *http.Request){
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Got json cred: %#v\n", creds)

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