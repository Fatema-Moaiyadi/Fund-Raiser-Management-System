package service

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
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
	Donate(request *models.DonationRequest) (*models.FundDetailsBrief, error)
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

	userDetails, err := fs.userDB.FindUser(constants.EmailColumnName, fundRequest.RaisedByUserEmail)
	if err != nil {
		return nil, err
	}

	fundDetails, err := fs.fundsDB.CreateFund(fundRequest, userDetails.UserID)
	if err != nil {
		return nil, err
	}

	return fundDetails, nil
}

func (fs *fundService) Donate(request *models.DonationRequest) (*models.FundDetailsBrief, error) {
	//get fund details
	fundDetails, err := fs.fundsDB.GetFundDetailsByID(request.DonatedInFund)
	if err != nil {
		return nil, err
	}

	totalRaisedAmount, err := fs.fundsDB.GetTotalRaisedAmountForFund(request.DonatedInFund)
	if err != nil {
		return nil, err
	}

	err = validations.ValidateDonateRequest(*request, totalRaisedAmount, fundDetails)
	if err != nil {
		return nil, err
	}

	existingDonations, err := fs.fundsDB.GetExistingDonationsForFundByUser(request.DonatedInFund, request.DonatedByUserID)
	if err != nil {
		return nil, err
	}

	if existingDonations == nil {
		//no active donations
		err = fs.fundsDB.CreateNewDonation(request)
	} else {
		err = fs.fundsDB.AddAmountToExistingDonation(request)
	}

	if err != nil {
		return nil, err
	}

	donatedFundDetails := &models.FundDetailsBrief{
		FundName:     fundDetails.FundName,
		AmountRaised: totalRaisedAmount,
		TotalAmount:  fundDetails.TotalAmount,
		FundStatus:   fundDetails.FundStatus,
	}

	return donatedFundDetails, nil
}
