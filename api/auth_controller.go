package api

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"log"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	DB *sql.DB
}

func (controller *AuthController) Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	uid := req.FormValue("uuid")
	password := req.FormValue("password")

	// authenticate device
	_, token, err := model.AuthDevice(controller.DB, uid, password)
	if err != nil {

		log.Println("Failed authentication attempt by ", req.RemoteAddr, req.UserAgent())
		http.Error(res, "Authentication failed", 401)
		return
	}

	// create json for token
	tokenJson, err := json.Marshal(&token)
	if err != nil {
		log.Println("Json: Error creating token: ", err)
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(tokenJson)
}