package database

const (
	getUserQuery               = "SELECT * from users where %s = $1 AND status = 'ACTIVE'"
	findUserForLoginQuery      = "SELECT * from users where email_id = $1 AND status = 'ACTIVE'"
	getUserInfoByIDQuery       = "SELECT user_id,email_id, first_name, last_name from users where %s AND status='ACTIVE'"
	insertUserQuery            = "INSERT INTO users (email_id,first_name, last_name,password,is_admin,status,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, 'ACTIVE', $6, $7) RETURNING user_id,email_id,first_name, last_name,password,is_admin,created_at,updated_at"
	updateUserSetClause        = "UPDATE users SET "
	whereUserIDClause          = ",updated_at=$%d where user_id=$%d"
	getUserByIDWithUpdatedInfo = "SELECT %s FROM users where user_id=$1"
	deleteUserByID             = "UPDATE users SET status='INACTIVE',updated_at=$1 where user_id=$2"
	getAllUserInfoQuery        = "SELECT user_id, email_id, first_name, last_name from users where status='ACTIVE'"

	insertFundQuery            = "INSERT INTO funds (raised_by_user_id,name,amount,status,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING fund_id,raised_by_user_id,name,amount,status,created_at,updated_at"
	getFundByIDQuery           = "SELECT * from funds where fund_id = $1"
	getActiveFundsRaisedByUser = "SELECT * from funds where raised_by_user_id = $1 AND status = 'IN_PROGRESS'"
	getAllActiveFunds          = "SELECT funds.name, funds.amount, sum(donation.amount) as amount_raised, users.email_id as raised_by from funds left join donation on funds.fund_id = donation.donated_in_fund_id inner join users on funds.raised_by_user_id = users.user_id where funds.status='IN_PROGRESS' group by funds.name,funds.amount,users.email_id"
	updateFundsSetClause       = "UPDATE funds SET "
	whereFundIDClause          = ",updated_at=$%d where fund_id=$%d AND status='IN_PROGRESS'"
	getFundByIDWithUpdateInfo  = "SELECT %s FROM funds where fund_id=$1"
	deleteFundByID             = "UPDATE funds SET status='DELETED', updated_at=$1 where fund_id=$2 AND status='IN_PROGRESS'"

	createDonationQuery                = "INSERT INTO donation (donated_in_fund_id, donated_by_user_id, amount, donation_status, created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	addAmountToExistingDonationQuery   = "UPDATE donation SET amount = amount+$1 where donated_in_fund_id=$2 AND donated_by_user_id=$3 AND donation_status = 'PAID'"
	getTotalRaisedAmountForFundQuery   = "SELECT sum(amount) from donation where donated_in_fund_id = $1"
	getPaidDonationsForFundByUserQuery = "SELECT * from donation where donated_in_fund_id=$1 AND donated_by_user_id=$2 AND donation_status = 'PAID'"
	getFundsRaisedByUserIDQuery        = "SELECT funds.name, funds.amount , sum(donation.amount) as amount_raised, funds.status from funds inner join donation on funds.fund_id = donation.donated_in_fund_id where funds.raised_by_user_id=$1 AND funds.status='IN_PROGRESS' group by funds.name,funds.amount,funds.status"
	getDonationsByUserIDQuery          = "SELECT funds.name, donation.amount as amount_donated, funds.status from funds inner join donation on funds.fund_id = donation.donated_in_fund_id where donation.donated_by_user_id=$1 AND donation.donation_status='PAID'"
	updateDonationStatusToRefund       = "UPDATE donation SET donation_status = 'REFUNDED', updated_at=$1 where donated_in_fund_id=$2"
)
