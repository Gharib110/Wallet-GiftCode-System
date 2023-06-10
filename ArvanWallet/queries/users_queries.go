package queries

const (
	InsertIntoUsersTable = `INSERT INTO users (name, phone_number, charge)
			VALUES ($1, $2, $3)`
	SelectUserByPhoneNumber = `SELECT id, name, phone_number, charge
		FROM public.users WHERE phone_number=$1`
	UpdateUserChargeByID = `UPDATE public.users SET charge=$1
	WHERE id=$2`
)
