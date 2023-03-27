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

	//initializing user flows
	userDb := database.NewUserDB(db)
	userService := service.NewUserService(userDb, jwtTokenService)
	userHandler := handler.NewUserHandler(userService)

	//initializing fund flows
	fundDB := database.NewFundsDB(db)
	fundService := service.NewFundService(fundDB, userDb)
	fundHandler := handler.NewFundsHandler(fundService, userService)

	deps := dependencies{
		userHandler:  userHandler,
		tokenService: jwtTokenService,
		fundsHandler: fundHandler,
	}
	return deps, nil
}
