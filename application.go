package main

import (
	"log"
	"net/http"
	"strconv"

	"os"

	"./conf"
	"./globals"
	"./handlers"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	defaultPort int = 5000
)

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
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	// log.SetOutput(f)

	conf.Setup()

	g := globals.Setup()
	defer globals.Close(g)

	http.Handle("/", route())

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}
	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
