package main

import (
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/config"
	"net/http"
)

func main() {

	err := config.ReadAndInitConfig("config.yml")
	if err != nil {
		fmt.Printf("Error in initialising config: %s", err.Error())
		return
	}

	deps, err := initDependencies()
	if err != nil {
		fmt.Printf("Error in initialising dependencies: %s", err.Error())
		return
	}

	fmt.Println("---------Starting Fund Raiser System---------")
	err = http.ListenAndServe(":"+config.GetAppPort(), getRoutes(deps))
	if err != nil {
		fmt.Printf("Error in starting app: %s", err.Error())
		return
	}
}
