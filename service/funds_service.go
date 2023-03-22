package service

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/database"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/fatema-moaiyadi/fund-raiser-system/validations"
)

type fundService struct {
	fundsDB database.FundsDB
	userDB  database.UserDatabase
}

type FundService interface {
	CreateFund(fundRequest *models.CreateFundRequest) (*models.FundDetails, error)
}

func NewFundService(fundsDB database.FundsDB, userDB database.UserDatabase) FundService {
	return &fundService{
		fundsDB: fundsDB,
		userDB:  userDB,
	}
}

func (fs *fundService) CreateFund(fundRequest *models.CreateFundRequest) (*models.FundDetails, error) {
	err := validations.ValidateCreateFundReq(fundRequest)
	if err != nil {
		return nil, err
	}

	userDetails, err := fs.userDB.FindUser(fundRequest.RaisedByUserEmail)
	if err != nil {
		return nil, err
	}

	fundDetails, err := fs.fundsDB.CreateFund(fundRequest, userDetails.UserID)
	if err != nil {
		return nil, err
	}

	return fundDetails, nil
}
