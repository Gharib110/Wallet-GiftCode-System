package queries2

const (
	InsertIntoRedemptionTable = `INSERT INTO public.user_redemptions(
	user_id, gift_code_id, type, redeemed_at)
	VALUES ($1, $2, $3, $4)`

	GetRedeemedRecords = `SELECT id, user_id, gift_code_id, type, redeemed_at
	FROM public.user_redemptions ORDER BY id DESC LIMIT 50`
)

