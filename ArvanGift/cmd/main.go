package main

import (
	"github.com/Gharib110/ArvanGift/api"
	models2 "github.com/Gharib110/ArvanGift/models"
	"github.com/Gharib110/ArvanGift/repo"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Dependency Injection and Initialization
func init() {
	app = appConfig{
		pQDB:        nil, // Bad Expose
		router:      nil,
		errorLogger: log.New(os.Stderr, "ERROR:\t", log.Ltime),
		infoLogger:  log.New(os.Stdout, "INFO:\t", log.Ltime),
	}

	pqDb := repo.NewGiftPQDB(app.errorLogger, app.infoLogger)
	app.pQDB = pqDb

	app.router = api.Routes()
	q := make(chan *models2.RedemptionQueueRequest)
	n := make(chan *models2.RedemptionQueueResponse)
	t := time.NewTicker(500 * time.Millisecond)

	ws := websocket.Upgrader{
		ReadBufferSize:  2048, // Size of the read buffer for incoming messages
		WriteBufferSize: 2048,
	}
	clients := make(map[string]*models2.WSConn)
	brc := make(chan *models2.UserRedemption)
	cApi := api.NewConfigAPI(app.pQDB, app.errorLogger, app.infoLogger,
		q, n, t, ws, clients, brc)
	app.cAPI = cApi
}

func main() {
	srv := &http.Server{
		Addr:              os.Getenv("GIFT_HOST") + ":" + os.Getenv("GIFT_PORT"),
		Handler:           app.router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       20 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC)

	go app.cAPI.RedemptionQueueHandler()
	go app.cAPI.BroadCasting()
	go func() {
		app.infoLogger.Println("HTTP1.x server is listening on " +
			os.Getenv("GIFT_HOST") + ":" + os.Getenv("GIFT_PORT"))
		if err := srv.ListenAndServe(); err != nil {
			app.errorLogger.Fatalln(err)
			return
		}
	}()

	<-sigC

	app.cAPI.CloseRedeemQueue()
	app.cAPI.CloseNotify()
	app.cAPI.CloseWSClients()
	app.cAPI.CloseBroadCaster()

	app.pQDB.DisposeDB()
}
