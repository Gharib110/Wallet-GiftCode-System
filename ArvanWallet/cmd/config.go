package main

import (
	"github.com/Gharib110/ArvanWallet/repo"
	"log"
	"net/http"
)

var app appConfig

type appConfig struct {
	pQDB        *repo.DBConfig
	router      http.Handler
	errorLogger *log.Logger
	infoLogger  *log.Logger
}
