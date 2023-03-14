package database

import (
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDBConnection() (*sqlx.DB, error) {
	dbConfig := config.GetDBConfig()
	dataSourceName := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable",
		dbConfig.UserName, dbConfig.Password, dbConfig.Port, dbConfig.DBName, dbConfig.Host)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		fmt.Printf("Error in connecting to database: %s", err.Error())
		fmt.Println()
		return nil, err
	}
	return db, nil
}
