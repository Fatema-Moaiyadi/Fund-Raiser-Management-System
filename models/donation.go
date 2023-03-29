package models

import "time"

type DonationRequest struct {
	DonationAmount   int64  `json:"donation_amount"`
	DonatedByEmailID string `json:"email_id"`
	DonatedByUserID  int64
	DonatedInFund    int64
}

//DonationAmountStatus enumeration for status of funds
type DonationAmountStatus int64

const (
	PAID DonationAmountStatus = iota
	REFUNDED
)

func (s DonationAmountStatus) String() string {
	switch s {
	case PAID:
		return "PAID"
	case REFUNDED:
		return "REFUNDED"
	}
	return "UNKNOWN"
}

type DonationData struct {
	DonationID           int64     `db:"donation_id"`
	DonationAmount       int64     `db:"amount"`
	DonatedInFundID      int64     `db:"donated_in_fund_id"`
	DonatedByUserID      int64     `db:"donated_by_user_id"`
	DonationAmountStatus string    `db:"donation_status"`
	CreatedAt            time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt            time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type FundDonationInfo struct {
	FundName            string `db:"name,omitempty" json:"fund_name,omitempty"`
	FundStatus          string `db:"status,omitempty" json:"fund_status,omitempty"`
	AmountDonatedByUser int64  `json:"amount_donated_by_user,omitempty"`
	TotalAmountRaised   int64  `json:"total_amount_raised,omitempty"`
	TotalAmount         int64  `db:"fund_amount,omitempty" json:"total_amount,omitempty"`
}

type DonationResponse struct {
	Code int `json:"code"`
	Data struct {
		FundInfo FundDonationInfo `json:"donation_info"`
	} `json:"data"`
}
