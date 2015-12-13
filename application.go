package main

import (
	"crypto/rand"
	"log"
	"net/http"

	"os"

	"./conf"
	"./handlers"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var redisDB int64 = 4

const redisSchemaVersion = 0
const (
	casUrl1 = "https://secure.its.yale.edu/cas/login?service="
	casUrl2 = "https://secure.its.yale.edu/cas/serviceValidate?"
)

func genRandomKey(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getCSRFKey() []byte {
	_authKey := os.Getenv("AUTH_KEY")
	if _authKey == "" {
		var err error
		authKey, err := genRandomKey(32)
		if err != nil {
			panic("Failed to generate random key for csrf.\n" + err.Error())
		}
		return authKey
	}

	return []byte(_authKey)
}

func route() http.Handler {
	r := mux.NewRouter().Schemes("http").Subrouter()

	r.HandleFunc("/", handlers.GetIndex).Methods("GET")
	r.HandleFunc("/urls/new", handlers.PostUrl).Methods("POST")
	r.HandleFunc("/urls/{key}", handlers.GetUrl).Methods("GET")
	r.HandleFunc("/urls/{key}/callback", hanlers.GetUrlCallback).Methods("GET")

	return csrf.Protect(getCSRFKey())(r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	conf := conf.Setup()
	defer credis.Close()
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	// log.SetOutput(f)

	http.Handle("/", route())

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
