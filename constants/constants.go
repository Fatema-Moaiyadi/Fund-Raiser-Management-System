package constants

import "github.com/lib/pq"

const (
	MinDonationAmount             = 100
	MaxDonationAmount             = 1000
	EmailColumnName               = "email_id"
	UserIDColumnName              = "user_id"
	ForeignKeyConstraintErrorCode = pq.ErrorCode("23503")
)
