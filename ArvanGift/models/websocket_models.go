package models2

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
)

type WsUser struct {
	UserName string `json:"username"`
}

type WSConn struct {
	Ws         *websocket.Conn
	OnlineData bool
	CloseChan  chan bool
	WriteChan  chan *UserRedemption
	CTX        context.Context
	CTXCancel  context.CancelFunc
	Username   string
	Lck        sync.Mutex
}
