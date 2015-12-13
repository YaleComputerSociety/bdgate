package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"strconv"

	"os"

	"./conf"
	"./handlers"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	defaultPort int = 5000
)

func genRandomKey(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func route() http.Handler {
	r := mux.NewRouter().Schemes("http").Subrouter()

	r.HandleFunc("/", handlers.GetIndex).Methods("GET")
	r.HandleFunc("/urls/new", handlers.PostUrl).Methods("POST")
	r.HandleFunc("/urls/{key}", handlers.GetUrl).Methods("GET")
	r.HandleFunc("/urls/{key}/callback", handlers.GetUrlCallback).Methods("GET")

	// Must be taken from the environment, because it must persist
	// accross instances and restarts.
	csrfKey := os.Getenv("CSRF_KEY_32")
	if csrfKey == "" {
		panic("CSRF_KEY_32 env variable not found.")
	}

	return csrf.Protect([]byte(csrfKey))(r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}

	config := conf.Setup()
	defer conf.Close(config)
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	// log.SetOutput(f)

	http.Handle("/", route())

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
