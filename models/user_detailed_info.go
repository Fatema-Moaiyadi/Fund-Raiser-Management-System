package models

type UserDetailedInfo struct {
	UserInfo        UserInfo           `json:"user_info"`
	RaisedFundsInfo []FundDetails      `json:"raised_funds_info"`
	DonationsInfo   []FundDonationInfo `json:"donations_info"`
}

type GetUserDetailsByIDResponse struct {
	Code int `json:"code"`
	Data struct {
		UserDetails UserDetailedInfo `json:"user_details"`
	} `json:"data"`
}
