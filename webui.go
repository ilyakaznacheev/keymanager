package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/new", NewKey)
	router.HandleFunc("/cancel/{key}", CancelKey).Methods("POST")
	router.HandleFunc("/status/{key}", CheckKey)
	router.HandleFunc("/info", Info)

	return router

}

func main() {

	initServer()
	defer shutDown()

	router := setupRouter()

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Panic(err)
	}

}
