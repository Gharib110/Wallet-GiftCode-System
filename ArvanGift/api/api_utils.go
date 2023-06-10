package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	models2 "github.com/Gharib110/ArvanGift/models"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

// dResponseWriter use for writing response to the user
func (c *ConfigAPI) dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		_, err := w.Write([]byte(data.(string)))
		return err
	} else if reflect.PtrTo(dataType).Kind() == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Struct == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Slice == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			c.errorLogger.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	}

	return errors.New("we could not be able to support data type that you passed")
}

// RedemptionQueueHandler Handle the queue should run in a separate go routine
func (c *ConfigAPI) RedemptionQueueHandler() {
	for d := range c.RedeemQueue {
		c.redemptionQueueRunnerHandler(d)
	}
}

// redemptionQueueRunnerHandler Process the requests in the queue
func (c *ConfigAPI) redemptionQueueRunnerHandler(d *models2.RedemptionQueueRequest) {
	req := &models2.RedeemTransactionDTO{
		UserID:     int64(d.U.ID),
		GiftCodeID: d.G.ID,
		Amount:     d.G.Amount,
		UserCharge: d.U.Charge,
		Type:       "Initial",
		RedeemedAt: time.Time{},
	}

	//TODO check the GiftCode Limitation
	if d.G.RedemptionCount == d.G.RedemptionLimit {
		resp := &models2.RedemptionQueueResponse{
			Identifier: d.I,
			UserRed: &models2.UserRedemption{
				UserID:     req.UserID,
				GiftCodeID: req.GiftCodeID,
				RedeemedAt: time.Now(),
				Type:       "FAILED",
			},
		}

		dao := &models2.UserRedemptionDAO{
			UserID:     req.UserID,
			GiftCodeID: req.GiftCodeID,
			RedeemedAt: time.Now(),
			Type:       "FAILED",
		}
		err := c.pQDB.CreateRedemption(d.R.Context(), dao)
		if err != nil {
			c.errorLogger.Println(err)
		}

		c.Notify <- resp
		return
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal the payload to JSON:", err)
	}

	// Send the POST request
	response, err := http.Post("http://localhost:9000/transaction/start",
		"application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println("Failed to make the request:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read the response body:", err)
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", response.StatusCode)
	}

	// Parse the JSON response

	rsp := &models2.RedeemTransactionDTO{
		UserID:     0,
		GiftCodeID: 0,
		Amount:     0,
		UserCharge: 0,
		Type:       "EMPTY",
		RedeemedAt: time.Time{},
	}
	err = json.Unmarshal(body, rsp)
	if err != nil {
		fmt.Println("Failed to parse the JSON response:", err)
	}

	if rsp.Type == "DONE" {
		//TODO Create User Redemption DONE record & increase the count of
		dao := &models2.UserRedemptionDAO{
			UserID:     rsp.UserID,
			GiftCodeID: rsp.GiftCodeID,
			RedeemedAt: rsp.RedeemedAt,
			Type:       rsp.Type,
		}
		err = c.pQDB.CreateRedemption(d.R.Context(), dao)
		if err != nil {
			c.errorLogger.Println(err)
			return
		}

		daoG := &models2.GiftCodeDAO{
			ID: d.G.ID,
		}
		err = c.pQDB.UpdateGiftCodeCount(d.R.Context(), daoG)
		if err != nil {
			c.errorLogger.Println(err)
			return
		}

		resp := &models2.RedemptionQueueResponse{
			Identifier: d.I,
			UserRed: &models2.UserRedemption{
				UserID:     rsp.UserID,
				GiftCodeID: rsp.GiftCodeID,
				RedeemedAt: rsp.RedeemedAt,
				Type:       rsp.Type,
			},
		}

		c.Notify <- resp
		return
	} else if rsp.Type == "FAILED" {
		//TODO Create User Redemption FAILED record
		dao := &models2.UserRedemptionDAO{
			UserID:     rsp.UserID,
			GiftCodeID: rsp.GiftCodeID,
			RedeemedAt: rsp.RedeemedAt,
			Type:       rsp.Type,
		}
		err = c.pQDB.CreateRedemption(d.R.Context(), dao)
		if err != nil {
			c.errorLogger.Println(err)
			return
		}

		resp := &models2.RedemptionQueueResponse{
			Identifier: d.I,
			UserRed: &models2.UserRedemption{
				UserID:     rsp.UserID,
				GiftCodeID: rsp.GiftCodeID,
				RedeemedAt: rsp.RedeemedAt,
				Type:       rsp.Type,
			},
		}

		c.Notify <- resp
		return
	}

	return
}

func (c *ConfigAPI) BroadCasting() {
	for redemption := range c.BroadCaster {
		for _, conn := range c.WSClients {
			if conn.OnlineData {
				conn.WriteChan <- redemption
			} else {
				continue
			}
		}
	}
}

func (c *ConfigAPI) CloseBroadCaster() {
	close(c.BroadCaster)
	c.infoLogger.Println("BroadCaster channel is closed")
}

func (c *ConfigAPI) CloseNotify() {
	close(c.Notify)
	c.infoLogger.Println("Notify channel is closed")
}

func (c *ConfigAPI) CloseRedeemQueue() {
	close(c.RedeemQueue)
	c.infoLogger.Println("RedemptionQueue channel is closed")
}

func (c *ConfigAPI) CloseWSClients() {
	for name, conn := range c.WSClients {
		close(conn.CloseChan)
		close(conn.WriteChan)
		err := conn.Ws.Close()
		if err != nil {
			delete(c.WSClients, name)
			continue
		}

		delete(c.WSClients, name)
	}

	c.infoLogger.Println("All WebSocket Connections and their associated channels are closed")
}
