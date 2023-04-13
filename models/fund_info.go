package models

import "time"

//FundStatus enumeration for status of funds
type FundStatus int

const (
	IN_PROGRESS FundStatus = iota
	DONE
	DELETED
)

func (s FundStatus) String() string {
	switch s {
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case DONE:
		return "DONE"
	case DELETED:
		return "DELETED"
	}
	return "UNKNOWN"
}

type CreateFundRequest struct {
	FundName          string `json:"fund_name,omitempty"`
	RaisedByUserEmail string `json:"raised_by,omitempty"`
	TotalAmount       int64  `json:"total_amount,omitempty"`
}

type FundDetails struct {
	FundID         int64      `db:"fund_id,omitempty" json:"fund_id,omitempty"`
	RaisedByUserID int64      `db:"raised_by_user_id,omitempty" json:"raised_by_user_id,omitempty"`
	FundName       string     `db:"name,omitempty" json:"fund_name,omitempty"`
	AmountRaised   int64      `db:"amount_raised,omitempty" json:"amount_raised,omitempty"`
	TotalAmount    int64      `db:"amount,omitempty" json:"total_amount,omitempty"`
	FundStatus     string     `db:"status,omitempty" json:"fund_status,omitempty"`
	CreatedAt      *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CreateFundResponse struct {
	Code int `json:"code"`
	Data struct {
		FundInfo FundDetails `json:"fund_info"`
	} `json:"data"`
}

type ActiveFundDetails struct {
	RaisedBy     string `db:"raised_by,omitempty" json:"raised_by_user_id,omitempty"`
	FundName     string `db:"name,omitempty" json:"fund_name,omitempty"`
	AmountRaised *int64 `db:"amount_raised,omitempty" json:"amount_raised,omitempty"`
	TotalAmount  int64  `db:"amount,omitempty" json:"total_amount,omitempty"`
}

type ActiveFundDetailsResponse struct {
	Code int `json:"code"`
	Data struct {
		FundsInfo []ActiveFundDetails `json:"funds_info"`
	} `json:"data"`
}

type UpdateFund struct {
	FundName        string `db:"name,omitempty" json:"fund_name,omitempty"`
	TotalFundAmount int64  `db:"amount,omitempty" json:"total_fund_amount,omitempty"`
	FundStatus      string `db:"status,omitempty"`
}

type UpdateFundResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UpdatedInfo UpdateFund `json:"updated_info"`
	} `json:"data"`
}

type DeleteFundResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
