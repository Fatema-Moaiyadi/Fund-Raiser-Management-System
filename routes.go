package main

import "github.com/gorilla/mux"

func getRoutes(deps dependencies) *mux.Router {
	r := mux.NewRouter()
	//initialise routes here

	r.HandleFunc("/login", deps.userHandler.LoginHandler()).Methods("POST")
	return r
}
