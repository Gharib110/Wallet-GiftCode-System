package api

import (
	"encoding/json"
	"github.com/Gharib110/ArvanWallet/models"
	"net/http"
)

type UsersHandlers interface {
	registerUser(http.ResponseWriter, *http.Request)
	getUserByPhone(http.ResponseWriter, *http.Request)
	checkAndRegister(http.ResponseWriter, *http.Request)
}

func (c *configAPI) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	// Default Init func
	var user = &models.User{
		Name:        "DEFAULT",
		PhoneNumber: "-1",
		Charge:      0,
	}

	//Decoder func
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	// Checking func
	if user.PhoneNumber == "-1" {
		http.Error(w, "Bad phone number !", http.StatusBadRequest)
		c.errorLogger.Println("Bad phone number !")
		return
	}

	userDao := &models.UserDAO{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Charge:      user.Charge,
	}

	err = c.pQDB.CreateUser(r.Context(), userDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.errorLogger.Println(err)
		return
	}

	err = c.dResponseWriter(w, user, http.StatusOK)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}

func (c *configAPI) getUserByPhone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	// Default Init func
	var user = &models.User{
		Name:        "DEFAULT",
		PhoneNumber: "-1",
		Charge:      0,
	}

	//Decoder func
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	// Checking func
	if user.PhoneNumber == "-1" {
		http.Error(w, "Bad phone number", http.StatusBadRequest)
		c.errorLogger.Println("Bad phone number")
		return
	}

	userDao := &models.UserDAO{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Charge:      user.Charge,
	}

	d := c.pQDB.GetUserInfoByPhoneNumber(r.Context(), userDao)
	if d == nil {
		c.errorLogger.Println("The data is nil!")
		http.Error(w, "The data is nil!", http.StatusInternalServerError)
		return
	}

	err = c.dResponseWriter(w, d, http.StatusOK)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}

// checkAndRegister check whether user exists or not and then register it
func (c *configAPI) checkAndRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		c.errorLogger.Println(r.RequestURI + r.Method + " is not available")
		return
	}

	// Default Init func
	var user = &models.User{
		Name:        "DEFAULT",
		PhoneNumber: "-1",
		Charge:      0,
	}

	//Decoder func
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.errorLogger.Println(err)
		return
	}

	// Checking func
	if user.PhoneNumber == "-1" {
		http.Error(w, "Bad Phone-Number !", http.StatusBadRequest)
		c.errorLogger.Println("Bad Phone-Number")
		return
	}

	userDao := &models.UserDAO{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Charge:      user.Charge,
	}

	// Handle the NewUserCreation
	d := c.pQDB.GetUserInfoByPhoneNumber(r.Context(), userDao)
	if d == nil {
		err = c.pQDB.CreateUser(r.Context(), userDao)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			c.errorLogger.Println(err)
			return
		}
		d = c.pQDB.GetUserInfoByPhoneNumber(r.Context(), userDao)
	}

	err = c.dResponseWriter(w, d, http.StatusOK)
	if err != nil {
		c.errorLogger.Println(err)
		return
	}
}
