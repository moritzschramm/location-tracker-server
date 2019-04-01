package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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

			log.Println("Unauthorized acces attempt (no token) by ", req.RemoteAddr, req.UserAgent(), err.Error())
			http.Error(res, "Unauthorized", 401)
			return
		}

		// check if token is valid
		token, err := model.GetAuthToken(controller.DB, tokenCookie.Value)
		if err != nil {

			log.Println("Unauthorized acces attempt (wrong token) by ", req.RemoteAddr, req.UserAgent(), err.Error())
			http.Error(res, "Forbidden", 403)
			return
		}

		// find corresponding device
		device, err := model.GetDevice(controller.DB, token.DeviceId)
		if err != nil {

			log.Println("Unauthorized access attempt by ", req.RemoteAddr, req.UserAgent(), err.Error())
			http.Error(res, "Forbidden", 403)
			return

		}

		// append device and token to request
		ctx := req.Context()
		ctx = context.WithValue(ctx, "device", device)
		ctx = context.WithValue(ctx, "token", token)

		next(res, req.WithContext(ctx))

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
		log.Println("Json: Error creating token: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	http.SetCookie(res, token.ToCookie())

	res.Write(tokenJson)
}

func (controller *AuthController) Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	token := req.Context().Value("token").(*model.AuthToken)

	http.SetCookie(res, token.UnsetCookie())
	token.Logout()

	res.WriteHeader(http.StatusOK)
}

func (controller *AuthController) TokenRefresh(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	token := req.Context().Value("token").(*model.AuthToken)

	token, err := token.Refresh()
	if err != nil {
		log.Println("Error refreshing token: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	tokenJson, err := json.Marshal(&token)
	if err != nil {
		log.Println("Json: Error creating token: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	http.SetCookie(res, token.ToCookie())

	res.Write(tokenJson)
}
