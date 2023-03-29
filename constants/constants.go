package constants

import "github.com/lib/pq"

const (
	MinDonationAmount             = 100
	MaxDonationAmount             = 1000
	EmailColumnName               = "email_id"
	PasswordColumnName            = "password"
	FirstNameColumnName           = "first_name"
	LastNameColumnName            = "last_name"
	UpdatedAtColumnName           = "updated_at"
	ForeignKeyConstraintErrorCode = pq.ErrorCode("23503")
)
