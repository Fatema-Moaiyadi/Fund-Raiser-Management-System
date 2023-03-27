package database

import (
	"database/sql"
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"time"
)

type fundsDB struct {
	database *sqlx.DB
}

type FundsDB interface {
	CreateFund(*models.CreateFundRequest, int64) (*models.FundDetails, error)
	CreateNewDonation(request *models.DonationRequest) error
	AddAmountToExistingDonation(request *models.DonationRequest) error
	GetTotalRaisedAmountForFund(fundID int64) (int64, error)
	GetFundDetailsByID(fundID int64) (*models.FundDetails, error)
	GetExistingDonationsForFundByUser(fundID, userID int64) ([]*models.DonationData, error)
}

func NewFundsDB(db *sqlx.DB) FundsDB {
	return &fundsDB{
		database: db,
	}
}

func (fdb *fundsDB) CreateFund(fundRequest *models.CreateFundRequest, userID int64) (*models.FundDetails, error) {
	fundDetails := new(models.FundDetails)

	err := fdb.database.Get(fundDetails, insertFundQuery,
		userID,
		fundRequest.FundName,
		fundRequest.TotalAmount,
		models.IN_PROGRESS.String(),
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}
	return fundDetails, nil
}

func (fdb *fundsDB) CreateNewDonation(request *models.DonationRequest) error {
	_, err := fdb.database.Exec(createDonationQuery,
		request.DonatedInFund,
		request.DonatedByUserID,
		request.DonationAmount,
		models.PAID.String(),
		time.Now(),
		time.Now())

	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == constants.ForeignKeyConstraintErrorCode {
				if strings.Contains(err.Error(), "fund_id") {
					return systemerrors.ErrFundNotFound
				} else if strings.Contains(err.Error(), "user_id") {
					return systemerrors.ErrUserNotFound
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (fdb *fundsDB) AddAmountToExistingDonation(request *models.DonationRequest) error {
	_, err := fdb.database.Exec(addAmountToExistingDonationQuery,
		request.DonationAmount,
		request.DonatedInFund,
		request.DonatedByUserID)

	if err != nil {
		return err
	}

	return nil
}

func (fdb *fundsDB) GetTotalRaisedAmountForFund(fundID int64) (int64, error) {
	var totalRaisedAmount interface{}
	err := fdb.database.Get(&totalRaisedAmount, getTotalRaisedAmountForFundQuery, fundID)
	if err != nil {
		return 0, err
	}

	if totalRaisedAmount == nil {
		return 0, nil
	}

	return totalRaisedAmount.(int64), nil
}

func (fdb *fundsDB) GetFundDetailsByID(fundID int64) (*models.FundDetails, error) {
	fundDetails := new(models.FundDetails)

	err := fdb.database.Get(fundDetails, getFundByIDQuery, fundID)
	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrFundNotFound
	}

	if err != nil {
		return nil, err
	}

	return fundDetails, nil
}

func (fdb *fundsDB) GetExistingDonationsForFundByUser(fundID, userID int64) ([]*models.DonationData, error) {
	var donationsByUserForFund []*models.DonationData
	err := fdb.database.Select(&donationsByUserForFund, getPaidDonationsForFundByUserQuery,
		fundID, userID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return donationsByUserForFund, nil
}
