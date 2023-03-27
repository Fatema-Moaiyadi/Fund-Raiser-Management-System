package database

const (
	findUserQuery   = "SELECT * from users where $1 = $2"
	insertUserQuery = "INSERT INTO users (email_id,name,password,is_admin,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id,email_id,name,password,is_admin,created_at,updated_at"

	insertFundQuery  = "INSERT INTO funds (raised_by_user_id,name,amount,status,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING fund_id,raised_by_user_id,name,amount,status,created_at,updated_at"
	getFundByIDQuery = "SELECT * from funds where fund_id = $1"

	createDonationQuery                = "INSERT INTO donation (donated_in_fund_id, donated_by_user_id, amount, donation_status, created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	addAmountToExistingDonationQuery   = "UPDATE donation SET amount = amount+$1 where donated_in_fund_id=$2 AND donated_by_user_id=$3 AND donation_status = 'PAID'"
	getTotalRaisedAmountForFundQuery   = "SELECT sum(amount) from donation where donated_in_fund_id = $1"
	getPaidDonationsForFundByUserQuery = "SELECT * from donation where donated_in_fund_id=$1 AND donated_by_user_id=$2 AND donation_status = 'PAID'"
)
