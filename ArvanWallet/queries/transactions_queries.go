package queries

const (
	InserIntoTransactionsTable = `INSERT INTO public.transactions(
	user_id, gift_code_id, amount, type, "timestamp")
	VALUES ($1, $2, $3, $4, $5)`
)
