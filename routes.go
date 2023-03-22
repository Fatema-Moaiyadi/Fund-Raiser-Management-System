package main

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/middleware"
	"github.com/gorilla/mux"
)

func getRoutes(deps dependencies) *mux.Router {
	r := mux.NewRouter()
	//initialise routes here

	r.HandleFunc("/login", deps.userHandler.LoginHandler()).Methods("POST")

	r.HandleFunc("/user", middleware.Authorize(deps.tokenService,
		deps.userHandler.CreateUser(), true)).Methods("POST")

	r.HandleFunc("/fund/create", middleware.Authorize(deps.tokenService,
		deps.fundsHandler.CreateFund(), true)).Methods("POST")
	return r
}
