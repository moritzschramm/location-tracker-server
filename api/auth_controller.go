package api

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"log"
	"context"
	"strings"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	DB *sql.DB
}

func (controller *AuthController) AuthenticationMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	if strings.HasPrefix(req.URL.Path, "/api") && req.URL.Path != "/api/login" {

		// get token from request
		tokenCookie, err := req.Cookie("token")
		if err != nil {

			log.Println("Unauthorized acces attempt (no token) by ", req.RemoteAddr, req.UserAgent())
			http.Error(res, "Unauthorized", 401)
		}

		// check if token is valid and find corresponding user
		device, err := model.CheckAuth(controller.DB, tokenCookie.Value)

		if err != nil {

			// append user to request
			ctx := req.Context()

			next(res, req.WithContext(context.WithValue(ctx, "device", device)))

		} else {

			log.Println("Unauthorized access attempt by ", req.RemoteAddr, req.UserAgent())
			http.Error(res, "Forbidden", 403)
		}

	} else {

		next(res, req)
	}
}

func (controller *AuthController) Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	uid := req.FormValue("uuid")
	password := req.FormValue("password")

	// authenticate device
	token, err := model.AuthDevice(controller.DB, uid, password)
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


func (controller *AuthController) Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {


}

func (controller *AuthController) TokenRefresh(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {


}