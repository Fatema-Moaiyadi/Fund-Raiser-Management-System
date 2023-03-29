package main

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/database"
	"github.com/fatema-moaiyadi/fund-raiser-system/handler"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
)

type dependencies struct {
	userHandler  handler.UserHandler
	tokenService service.TokenService
	fundsHandler handler.FundsHandler
}

func initDependencies() (dependencies, error) {
	db, err := database.InitDBConnection()
	if err != nil {
		return dependencies{}, err
	}

	jwtTokenService := service.NewJWTTokenService()

	//initialising databases
	fundDB := database.NewFundsDB(db)
	userDb := database.NewUserDB(db)

	//initializing services
	userService := service.NewUserService(userDb, jwtTokenService, fundDB)
	fundService := service.NewFundService(fundDB, userDb)

	//initialising handlers
	userHandler := handler.NewUserHandler(userService)
	fundHandler := handler.NewFundsHandler(fundService, userService)

	deps := dependencies{
		userHandler:  userHandler,
		tokenService: jwtTokenService,
		fundsHandler: fundHandler,
	}
	return deps, nil
}
