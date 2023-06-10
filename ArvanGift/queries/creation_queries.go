package queries2

const (
	CreateGiftCodesTable = `CREATE TABLE IF NOT EXISTS gift_codes (
  id SERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL UNIQUE,
  amount NUMERIC(10, 2) NOT NULL,
  is_active BOOLEAN NOT NULL,
  redemption_limit INT NOT NULL,
  redemption_count INT NOT NULL,
  start_time TIMESTAMP NOT NULL,
  expiration_time TIMESTAMP NOT NULL
)`

	CreateRedemptionTable = `CREATE TABLE IF NOT EXISTS user_redemptions (
  id SERIAL PRIMARY KEY,
  user_id INT,
  gift_code_id INT REFERENCES gift_codes(id),
  type VARCHAR(255) NOT NULL,
  redeemed_at TIMESTAMP NOT NULL DEFAULT NOW()
)`
)
