package models

import "time"

//FundStatus enumeration for status of funds
type FundStatus int64

const (
	IN_PROGRESS FundStatus = iota
	DONE
	APPLIED
	DELETED
)

func (s FundStatus) String() string {
	switch s {
	case IN_PROGRESS:
		return "IN_PROGRESS"
	case DONE:
		return "DONE"
	case APPLIED:
		return "APPLIED"
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
	AmountRaised   int64      `json:"amount_raised,omitempty"`
	TotalAmount    int64      `db:"amount,omitempty" json:"total_amount,omitempty"`
	FundStatus     FundStatus `db:"status,omitempty" json:"fund_status,omitempty"`
	CreatedAt      time.Time  `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time  `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type FundDetailsBrief struct {
	FundName     string     `db:"name,omitempty" json:"fund_name,omitempty"`
	AmountRaised int64      `json:"amount_raised,omitempty"`
	TotalAmount  int64      `db:"amount,omitempty" json:"total_amount,omitempty"`
	FundStatus   FundStatus `db:"status,omitempty" json:"fund_status,omitempty"`
}

type CreateFundResponse struct {
	Code int `json:"code"`
	Data struct {
		FundInfo FundDetails `json:"fund_info"`
	} `json:"data"`
}


