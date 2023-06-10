package api

import (
	"github.com/Gharib110/ArvanWallet/repo"
	"log"
)

type configAPI struct {
	pQDB        *repo.DBConfig
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

var conf configAPI

func NewConfigAPI(pqDb *repo.DBConfig, errLog *log.Logger, infoLog *log.Logger) {
	conf = configAPI{
		pQDB:        pqDb,
		errorLogger: errLog,
		infoLogger:  infoLog,
	}
}
