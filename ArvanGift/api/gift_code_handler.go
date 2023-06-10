package api

import (
	"encoding/json"
	"github.com/Gharib110/ArvanGift/models"
	"net/http"
	"time"
)

type GiftCodeHandlers interface {
	registerGiftCode(http.ResponseWriter, *http.Request)
	getGiftCodeByCode(http.ResponseWriter, *http.Request)
}

func (c *ConfigAPI) registerGiftCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	// Default Init func
	var giftCode = &models2.GiftCode{
		Code:            "DEFAULT",
		Amount:          -1,
		IsActive:        false,
		RedemptionLimit: -1,
		RedemptionCount: 0,
		StartTime:       "DEFAULT",
		ExpirationTime:  "DEFAULT",
	}

	//Decoder func
	err := json.NewDecoder(r.Body).Decode(&giftCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	st, err := time.Parse("2006-01-02 15:04:05", giftCode.StartTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println("Bad Information !")
		return
	}
	exp, err := time.Parse("2006-01-02 15:04:05", giftCode.ExpirationTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println("Bad Information !")
		return
	}

	// Checking func
	if giftCode.Code == "DEFAULT" || giftCode.Amount == -1 ||
		giftCode.RedemptionLimit == -1 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println("Bad Information !")
		return
	}

	giftDao := &models2.GiftCodeDAO{
		Code:            giftCode.Code,
		Amount:          giftCode.Amount,
		IsActive:        giftCode.IsActive,
		RedemptionLimit: giftCode.RedemptionLimit,
		RedemptionCount: giftCode.RedemptionCount,
		StartTime:       st,
		ExpirationTime:  exp,
	}

	err = c.pQDB.RegisterGiftCode(r.Context(), giftDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.errorLogger.Println(err)
		return
	}

	err = c.dResponseWriter(w, giftDao, http.StatusOK)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}

func (c *ConfigAPI) getGiftCodeByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	var giftCode = models2.GiftCode{
		Code:            "DEFAULT",
		Amount:          -1,
		IsActive:        false,
		RedemptionLimit: -1,
		RedemptionCount: 0,
		StartTime:       "DEFAULT",
		ExpirationTime:  "DEFAULT",
	}

	err := json.NewDecoder(r.Body).Decode(&giftCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	giftDao := models2.GiftCodeDAO{
		Code: giftCode.Code,
	}

	g := c.pQDB.GetGiftCodeByID(r.Context(), &giftDao)
	if g.ID == 0 {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		c.errorLogger.Println("Internal Error")
		return
	}

	err = c.dResponseWriter(w, g, http.StatusOK)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}
