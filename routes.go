package main

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/middleware"
	"github.com/gorilla/mux"
)

func getRoutes(deps dependencies) *mux.Router {
	r := mux.NewRouter()
	//initialise routes here

	r.HandleFunc("/login", deps.userHandler.LoginHandler()).Methods("POST")

	r.HandleFunc("/user", middleware.AuthorizeAdmin(deps.tokenService,
		deps.userHandler.CreateUser())).Methods("POST")

	r.HandleFunc("/fund/create", middleware.AuthorizeAdmin(deps.tokenService,
		deps.fundsHandler.CreateFund())).Methods("POST")

	r.HandleFunc("/donate/{fund_id}", middleware.Authorize(deps.tokenService,
		deps.fundsHandler.DonateInFund())).Methods("POST")

	r.HandleFunc("/user/update", middleware.Authorize(deps.tokenService,
		deps.userHandler.UpdateUserInfo())).Methods("PUT")

	r.HandleFunc("/user", middleware.AuthorizeAdmin(deps.tokenService,
		deps.userHandler.DeleteUserByID())).Methods("DELETE")

	r.HandleFunc("/user/{user_id}", middleware.Authorize(deps.tokenService,
		deps.userHandler.GetUserInfoByID())).Methods("GET")
	return r
}
