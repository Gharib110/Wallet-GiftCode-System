package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	models2 "github.com/Gharib110/ArvanGift/models"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type RedemptionHandlers interface {
	redeem(http.ResponseWriter, *http.Request)
	getRedeemedRecords(http.ResponseWriter, *http.Request)
}

type RedemptionUtils interface {
	checkGiftCodeExistence(string, *http.Request) (*models2.GiftCode, error)
	checkUserExistence(string) (*models2.UserInfoDTO, error)
}

type WebsocketUtils interface {
	handleWebsocketConnection(*models2.WSConn)
}

func (c *ConfigAPI) redeem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	var info = &models2.UserRedemptionInfo{
		PhoneNumber: "-1",
		GiftCode:    "-1",
	}

	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	if info.PhoneNumber == "-1" || info.GiftCode == "-1" {
		http.Error(w, "Bad information provided", http.StatusBadRequest)
		c.errorLogger.Println("Bad information provided")
		return
	}

	// Check the Existence of GiftCode
	giftCode, err := c.checkGiftCodeExistence(info.GiftCode, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.errorLogger.Println(err)
		return
	}

	exp, _ := time.Parse("2006-01-02 15:04:05", giftCode.ExpirationTime)
	if time.Now().After(exp) {
		http.Error(w, "The GIFT CODE is Expired", http.StatusBadRequest)
		c.errorLogger.Println("The GIFT CODE is Expired")
		return
	}

	// Check the Existence of User
	user, err := c.checkUserExistence(info.PhoneNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.errorLogger.Println(err)
		return
	}

	dao := &models2.UserRedemptionDAO{
		UserID:     int64(user.ID),
		GiftCodeID: giftCode.ID,
		RedeemedAt: time.Now(),
		Type:       "CREATED",
	}

	err = c.pQDB.CreateRedemption(r.Context(), dao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.errorLogger.Println(err)
		return
	}

	i := user.PhoneNumber + giftCode.Code + time.Now().Format("2006-01-02 15:04:05.99999999")
	redeemD := &models2.RedemptionQueueRequest{
		R: r,
		W: w,
		U: user,
		G: giftCode,
		I: i,
	}

	c.RedeemQueue <- redeemD

	for {
		select {
		case msg := <-c.Notify:
			if msg.Identifier == redeemD.I {
				c.BroadCaster <- msg.UserRed
				c.dResponseWriter(w, msg.UserRed, http.StatusOK)
				return
			}
		case <-c.T.C:
			continue
		}
	}
}

func (c *ConfigAPI) checkGiftCodeExistence(giftCode string,
	r *http.Request) (*models2.GiftCode, error) {
	gDao := &models2.GiftCodeDAO{
		ID:              0,
		Code:            giftCode,
		Amount:          0,
		IsActive:        false,
		RedemptionLimit: 0,
		RedemptionCount: 0,
		StartTime:       time.Time{},
		ExpirationTime:  time.Time{},
	}

	gift := c.pQDB.GetGiftCodeByID(r.Context(), gDao)
	if gift.ID == 0 {
		c.errorLogger.Println("could not find the specific gift-code")
		return nil, errors.New("could not find the specific gift-code")
	}

	return gift, nil
}

func (c *ConfigAPI) checkUserExistence(pn string) (*models2.UserInfoDTO, error) {
	payload := map[string]interface{}{
		"phone_number": pn,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal the payload to JSON:", err)
		return nil, err
	}

	// Send the POST request
	response, err := http.Post("http://localhost:9000/user/check-register",
		"application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to make the request:", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read the response body:", err)
		return nil, err
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", response.StatusCode)
		return nil, err
	}

	// Parse the JSON response
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Failed to parse the JSON response:", err)
		return nil, err
	}

	fmt.Println("Request successful!")
	fmt.Println("Response Data:", responseData)

	user := &models2.UserInfoDTO{
		ID:          responseData["id"].(float64),
		Name:        responseData["name"].(string),
		PhoneNumber: responseData["phone_number"].(string),
		Charge:      responseData["charge"].(float64),
	}
	return user, nil
}

func (c *ConfigAPI) getRedeemedRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}
	conn, err := c.WsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	ch := make(chan bool)
	wCh := make(chan *models2.UserRedemption)

	info := &models2.WsUser{UserName: "User;" + time.Now().String()}
	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(`TCP/HTTP connection is upgraded to TCP/WEBSOCKET connection Dear `+
			info.UserName+`Commands: SHOW-DESCRIPTION, CLOSE, ENABLE-ONLINE, DISABLE-ONLINE`))
	if err != nil {
		c.errorLogger.Println(fmt.Sprintf("%v", r))
		c.dResponseWriter(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	newC := &models2.WSConn{
		Ws:         conn,
		OnlineData: false,
		CloseChan:  ch,
		WriteChan:  wCh,
		CTX:        ctx,
		CTXCancel:  cancel,
		Username:   info.UserName,
		Lck:        sync.Mutex{},
	}

	c.WSClients[newC.Username] = newC
	go c.handleWebsocketConnection(newC)
}

func (c *ConfigAPI) handleWebsocketConnection(ws *models2.WSConn) {
	defer func() {
		if r := recover(); r != nil {
			c.errorLogger.Println(fmt.Sprintf("%v", r))
		}
	}()

	// Listening WS Routine
	go func() {
		for {
			msgType, data, err := ws.Ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.errorLogger.Println("Normal Closure: " + err.Error())
					ws.CloseChan <- true
					return
				} else if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
					c.errorLogger.Println("Abnormal Closure: " + err.Error())
					ws.CloseChan <- true
					return
				} else {
					c.errorLogger.Println(err)
					ws.CloseChan <- true
					return
				}
			}

			if msgType != websocket.TextMessage {
				err = ws.Ws.WriteMessage(websocket.TextMessage, []byte("Bad Data !, Should be in string"))
				if err != nil {
					c.errorLogger.Println(err)
					continue
				}
			}

			if string(data) == "SHOW-REDEMPTION" {
				records, err := c.pQDB.GetRedeemedRecords(ws.CTX)
				if err != nil {
					c.errorLogger.Println(err)
					continue
				}

				length := records.Len()
				for i := 0; i < length; i++ {
					ws.WriteChan <- records.Back().Value.(*models2.UserRedemption)
					records.Remove(records.Back())
				}
				ws.Lck.Lock()
				err = ws.Ws.WriteMessage(websocket.TextMessage, []byte("SHOW-REDEMPTION is FINISHED"))
				ws.Lck.Unlock()
				if err != nil {
					c.errorLogger.Println(err)
					continue
				}
				continue
			} else if string(data) == "CLOSE" {
				ws.CloseChan <- true
				return
			} else if string(data) == "ENABLE-ONLINE" {
				ws.OnlineData = true
				ws.Lck.Lock()
				err = ws.Ws.WriteMessage(websocket.TextMessage,
					[]byte("Online-Data Reading is set to TRUE"))
				ws.Lck.Unlock()
			} else if string(data) == "DISABLE-ONLINE" {
				ws.OnlineData = false
				ws.Lck.Lock()
				err = ws.Ws.WriteMessage(websocket.TextMessage,
					[]byte("Online-Data Reading is set to FALSE"))
				ws.Lck.Unlock()
			}
		}
	}()

	// Speaking WS Routine
	go func() {
		for redemption := range ws.WriteChan {
			ws.Lck.Lock()
			err := ws.Ws.WriteJSON(redemption)
			ws.Lck.Unlock()
			if err != nil {
				c.errorLogger.Println(err)
				continue
			}
		}
	}()

	<-ws.CloseChan
	ws.Lck.Lock()
	err := ws.Ws.WriteMessage(websocket.TextMessage, []byte("Bye, Bye "+ws.Username))
	ws.Lck.Unlock()

	if err != nil {
		c.errorLogger.Println(err)
	}

	err = ws.Ws.Close()
	ws.CTXCancel()
	close(ws.WriteChan)
	close(ws.CloseChan)
	delete(c.WSClients, ws.Username)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}
