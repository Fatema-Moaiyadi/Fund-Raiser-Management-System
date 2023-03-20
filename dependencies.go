package main

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/database"
	"github.com/fatema-moaiyadi/fund-raiser-system/handler"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
)

type dependencies struct {
	userHandler handler.UserHandler
}

func initDependencies() (dependencies, error) {
	db, err := database.InitDBConnection()
	if err != nil {
		return dependencies{}, err
	}

	jwtTokenService := service.NewJWTTokenService()

	userDb := database.NewUserDB(db)
	userService := service.NewUserService(userDb, jwtTokenService)
	userHandler := handler.NewUserHandler(userService)

	deps := dependencies{
		userHandler: userHandler,
	}
	return deps, nil
}
