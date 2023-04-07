package database

import (
	"database/sql"
	"fmt"
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
	CreateNewDonation(request *models.DonationRequest, totalRaisedAmount int64, totalFundAmount int64) error
	AddAmountToExistingDonation(request *models.DonationRequest, totalRaisedAmount int64, totalFundAmount int64) error
	GetTotalRaisedAmountForFund(fundID int64) (int64, error)
	GetFundDetailsByID(fundID int64) (*models.FundDetails, error)
	GetExistingDonationsForFundByUser(fundID, userID int64) ([]*models.DonationData, error)
	GetFundsRaisedByUserID(userID int64) ([]*models.FundDetails, error)
	GetAllActiveFunds() ([]models.ActiveFundDetails, error)
	UpdateFundByID(updateParams map[string]interface{}, fundID int64) (*models.UpdateFund, error)
	DeleteFundByID(fundID int64) error
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

func (fdb *fundsDB) CreateNewDonation(request *models.DonationRequest, totalRaisedAmount int64, totalFundAmount int64) error {
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

	if totalRaisedAmount+request.DonationAmount == totalFundAmount {
		updateParams := make(map[string]interface{})
		updateParams[constants.FundStatusColumnName] = models.DONE.String()
		_, err = fdb.UpdateFundByID(updateParams, request.DonatedInFund)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fdb *fundsDB) AddAmountToExistingDonation(request *models.DonationRequest, totalRaisedAmount int64, totalFundAmount int64) error {
	_, err := fdb.database.Exec(addAmountToExistingDonationQuery,
		request.DonationAmount,
		request.DonatedInFund,
		request.DonatedByUserID)

	if err != nil {
		return err
	}

	if totalRaisedAmount+request.DonationAmount == totalFundAmount {
		updateParams := make(map[string]interface{})
		updateParams[constants.FundStatusColumnName] = models.DONE.String()
		_, err = fdb.UpdateFundByID(updateParams, request.DonatedInFund)
		if err != nil {
			return err
		}
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

func (fdb *fundsDB) GetFundsRaisedByUserID(userID int64) ([]*models.FundDetails, error) {
	var raisedFunds []*models.FundDetails

	err := fdb.database.Select(&raisedFunds, getActiveFundsRaisedByUser, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return raisedFunds, nil
}

func (fdb *fundsDB) GetAllActiveFunds() ([]models.ActiveFundDetails, error) {
	var activeFunds []models.ActiveFundDetails

	err := fdb.database.Select(&activeFunds, getAllActiveFunds)
	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrNoActiveFunds
	}

	if err != nil {
		return nil, err
	}

	return activeFunds, nil
}

func (fdb *fundsDB) UpdateFundByID(updateParams map[string]interface{}, fundID int64) (*models.UpdateFund, error) {
	var (
		setValues []interface{}
		setKeys   []string
		getValues []string
		index     = 1
	)

	for key, value := range updateParams {
		setKeys = append(setKeys, key+"=$"+fmt.Sprintf("%d", index))
		setValues = append(setValues, value)
		getValues = append(getValues, key)
		index++
	}

	updateQuery := updateFundsSetClause + strings.Join(setKeys, ",") + fmt.Sprintf(whereFundIDClause, index, index+1)

	tx, err := fdb.database.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(updateQuery, append(setValues, time.Now(), fundID)...)

	affectedRows, e := result.RowsAffected()
	if e != nil {
		return nil, e
	}

	if affectedRows == 0 {
		return nil, systemerrors.ErrFundNotFound
	}
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return nil, e
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedInfo := new(models.UpdateFund)
	getQuery := fmt.Sprintf(getFundByIDWithUpdateInfo, strings.Join(getValues, ","))
	err = fdb.database.Get(updatedInfo, getQuery, fundID)

	if err != nil {
		return nil, err
	}

	return updatedInfo, nil
}

func (fdb *fundsDB) DeleteFundByID(fundID int64) error {
	tx, err := fdb.database.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(deleteFundByID, time.Now(), fundID)
	affectedRows, e := result.RowsAffected()
	if e != nil {
		return e
	}

	if affectedRows == 0 {
		return systemerrors.ErrFundNotFound
	}
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}

	//if fund will be deleted, all donations should be refunded
	_, err = tx.Exec(updateDonationStatusToRefund, time.Now(), fundID)
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
