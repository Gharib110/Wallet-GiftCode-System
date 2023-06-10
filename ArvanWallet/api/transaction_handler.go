package api

import (
	"encoding/json"
	"github.com/Gharib110/ArvanWallet/models"
	"net/http"
	"time"
)

type TransactionHandler interface {
	startTransaction(http.ResponseWriter, *http.Request)
}

func (c *configAPI) startTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	trData := &models.TransactionDTO{
		UserID:     0,
		GiftCodeID: 0,
		Amount:     0,
		UserCharge: 0,
		Type:       "",
		Timestamp:  time.Time{},
	}

	err := json.NewDecoder(r.Body).Decode(trData)
	if err != nil {
		http.Error(w, "Bad Data Provided !", http.StatusBadRequest)
		c.errorLogger.Println("Bad Data Provided !")
		return
	}

	err, dto := c.pQDB.BeginTransaction(r.Context(), trData)
	if err != nil {
		return
	}

	err = c.dResponseWriter(w, dto, http.StatusOK)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		c.errorLogger.Println("Internal Error")
		return
	}

	return
}
