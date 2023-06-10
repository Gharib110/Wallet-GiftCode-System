package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	dB          *sql.DB
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

var dbConfig *DBConfig

func NewGiftPQDB(errLog *log.Logger, infoLog *log.Logger) *DBConfig {
	dbConfig = &DBConfig{
		dB:          nil,
		errorLogger: errLog,
		infoLogger:  infoLog,
	}

	port, _ := strconv.Atoi(os.Getenv("GIFT_DB_PORT"))

	dbname := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("GIFT_DB_HOST"), port,
		os.Getenv("GIFT_DB_USER"),
		os.Getenv("GIFT_DB_PASSWORD"), os.Getenv("GIFT_DB_NAME"))

	giftDb, err := sql.Open("postgres", dbname)
	if err != nil {
		dbConfig.errorLogger.Fatalln(err)
	} else {
		dbConfig.infoLogger.Println("DB is initiated.")
	}

	dbConfig.dB = giftDb
	dbConfig.createTables()

	return dbConfig
}

// PingPQDB pings database
func PingPQDB() {
	err := dbConfig.dB.Ping()
	if err != nil {
		dbConfig.errorLogger.Panicln(err)
	} else {
		dbConfig.infoLogger.Println("DB is Online")
	}
}

func (c *DBConfig) DisposeDB() {
	c.deleteTables()
	err := c.dB.Close()
	if err != nil {
		return
	}
}
