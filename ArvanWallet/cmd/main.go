package main

import (
	"github.com/Gharib110/ArvanWallet/api"
	"github.com/Gharib110/ArvanWallet/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	app = appConfig{
		pQDB:        nil,
		router:      nil,
		errorLogger: log.New(os.Stderr, "ERROR:\t", log.Ltime),
		infoLogger:  log.New(os.Stdout, "INFO:\t", log.Ltime),
	}

	pqDb := repo.NewWalletPQDB(app.errorLogger, app.infoLogger)
	app.pQDB = pqDb

	app.router = api.Routes()
	api.NewConfigAPI(app.pQDB, app.errorLogger, app.infoLogger)
}

func main() {

	srv := &http.Server{
		Addr:              os.Getenv("WALLET_HOST") + ":" + os.Getenv("WALLET_PORT"),
		Handler:           app.router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC)

	go func() {
		app.infoLogger.Println("HTTP1.x server is listening on " +
			os.Getenv("WALLET_HOST") + ":" + os.Getenv("WALLET_PORT"))
		if err := srv.ListenAndServe(); err != nil {
			app.errorLogger.Fatalln(err)
			return
		}
	}()

	<-sigC

	app.pQDB.DisposeDB()
}
