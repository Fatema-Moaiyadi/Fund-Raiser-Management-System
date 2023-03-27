package models

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
	DonationID           int64                `db:"donation_id"`
	DonationAmount       int64                `db:"donation_amount"`
	DonatedInFundID      int64                `db:"donated_in_fund_id"`
	DonatedByUserID      int64                `db:"donated_by_user_id"`
	DonationAmountStatus DonationAmountStatus `db:"donation_status"`
}

type DonationResponse struct {
	Code int `json:"code"`
	Data struct {
		FundInfo FundDetailsBrief `json:"donation_info"`
	} `json:"data"`
}
