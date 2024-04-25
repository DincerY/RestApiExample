package controllers

import (
	"RestApiExample/models"
	u "RestApiExample/utils"
	"github.com/gorilla/mux"
	"net/http"
)

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["userId"]

	resp := models.GetUser(id)
	resp["func"] = "GetUser"
	u.Respond(w, resp)
}
