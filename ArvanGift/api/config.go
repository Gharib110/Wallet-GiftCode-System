package api

import (
	"github.com/Gharib110/ArvanGift/models"
	"github.com/Gharib110/ArvanGift/repo"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ConfigAPI struct {
	pQDB        *repo.DBConfig
	RedeemQueue chan *models2.RedemptionQueueRequest
	Notify      chan *models2.RedemptionQueueResponse
	T           *time.Ticker
	WsUpgrade   websocket.Upgrader
	WSClients   map[string]*models2.WSConn
	BroadCaster chan *models2.UserRedemption
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

var conf ConfigAPI

func NewConfigAPI(pqDb *repo.DBConfig, errLog *log.Logger, infoLog *log.Logger,
	q chan *models2.RedemptionQueueRequest,
	n chan *models2.RedemptionQueueResponse,
	t *time.Ticker, ws websocket.Upgrader,
	clients map[string]*models2.WSConn,
	brc chan *models2.UserRedemption) *ConfigAPI {
	conf = ConfigAPI{
		pQDB:        pqDb,
		errorLogger: errLog,
		infoLogger:  infoLog,
		RedeemQueue: q,
		Notify:      n,
		T:           t,
		WsUpgrade:   ws,
		WSClients:   clients,
		BroadCaster: brc,
	}

	return &conf
}
