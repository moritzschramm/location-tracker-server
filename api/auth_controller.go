package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

// authentication middleware checks if all calls to the API (except login) are authenticated
// static files are always served
func (controller *Controller) AuthenticationMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

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

// login handler, if correct UUID and password are entered,
// user gets logged in (receives token, also set via HTTP cookie)
func (controller *Controller) Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

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

// logout handler, removes cookie and token from database
func (controller *Controller) Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	token := req.Context().Value("token").(*model.AuthToken)

	http.SetCookie(res, token.UnsetCookie())
	token.Logout()

	res.WriteHeader(http.StatusOK)
}

// refreshes old (but not expired) token and returns new token to user
// old token is not usable after refresh
func (controller *Controller) TokenRefresh(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

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


// checks if request was performed by admin user (device)
// automatically logs unauthorized access and generates HTTP response
func (controller *Controller) CheckIfAdmin(req *http.Request) bool {

	user := req.Context().Value("device").(*model.Device)
	if user.UUID.String() != controller.Config.AdminUUID {
		log.Println("Unauthorized attempt to create new device by ", req.RemoteAddr, req.UserAgent())
		http.Error(res, "Unauthorized", 403)
		return false
	}
	return true
}