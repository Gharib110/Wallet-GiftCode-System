package queries2

const (
	InsertIntoGiftCodeTable = `INSERT INTO public.gift_codes(
	code, amount, is_active, redemption_limit, redemption_count, start_time, expiration_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	GetGiftCodeByCodeID = `SELECT id, code, amount, is_active, 
       redemption_limit, redemption_count, start_time, expiration_time
	FROM public.gift_codes WHERE code=$1`

	UpdateGiftCodeRedemptionCount = `UPDATE public.gift_codes
	SET redemption_count=redemption_count+1 
	WHERE id=$1`
)
