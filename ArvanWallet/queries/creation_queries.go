package queries

const (
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone_number VARCHAR(20) NOT NULL UNIQUE,
  charge BIGINT NOT NULL DEFAULT 0
)`
	CreateTransactionsTable = `CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  gift_code_id INT,
  amount NUMERIC(10, 2) NOT NULL,
  type VARCHAR(255) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT NOW()
)`
)
