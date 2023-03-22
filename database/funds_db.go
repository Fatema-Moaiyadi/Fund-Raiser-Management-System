package database

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type fundsDB struct {
	database *sqlx.DB
}

type FundsDB interface {
	CreateFund(*models.CreateFundRequest, int) (*models.FundDetails, error)
}

func NewFundsDB(db *sqlx.DB) FundsDB {
	return &fundsDB{
		database: db,
	}
}

func (fdb *fundsDB) CreateFund(fundRequest *models.CreateFundRequest, userID int) (*models.FundDetails, error) {
	fundDetails := new(models.FundDetails)

	err := fdb.database.Get(fundDetails, insertFundQuery,
		userID,
		fundRequest.FundName,
		fundRequest.TotalAmount,
		models.IN_PROGRESS,
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}
	return fundDetails, nil
}
