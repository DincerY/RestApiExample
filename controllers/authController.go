package controllers

import (
	"RestApiExample/models"
	u "RestApiExample/utils"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := account.Create()
	resp["func"] = "CreateAccount"
	u.Respond(w, resp)

}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	resp := models.Login(w, account.Email, account.Password)

	resp["func"] = "Authenticate"
	u.Respond(w, resp)
}
