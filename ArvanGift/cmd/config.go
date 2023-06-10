package main

import (
	"github.com/Gharib110/ArvanGift/api"
	"github.com/Gharib110/ArvanGift/repo"
	"log"
	"net/http"
)

var app appConfig

type appConfig struct {
	pQDB        *repo.DBConfig
	cAPI        *api.ConfigAPI
	router      http.Handler
	errorLogger *log.Logger
	infoLogger  *log.Logger
}
