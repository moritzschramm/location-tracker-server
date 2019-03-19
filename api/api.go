package api

import (
	"database/sql"
	"net/http"
	"strings"
	"log"

	"github.com/moritzschramm/location-tracker-server/config"
	"github.com/moritzschramm/location-tracker-server/mqtt"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func SetupAPI(db *sql.DB, mqttClient MQTT.Client, config config.Config) {

	// setup router
	router := httprouter.New()
	//router.NotFound = http.HandlerFunc(NotFoundHandler)

	// static files
	//router.GET("/", Index)
	router.ServeFiles("/assets/*filepath", http.Dir(config.PublicDir+"/assets"))

	// api
	deviceController := &DeviceController{DB: db, Mqtt: mqtt.User{config.MQTT}}
	router.POST("/api/device/new", 			deviceController.NewDevice)
	router.POST("/api/device/delete/:uid", 	deviceController.DeleteDevice)

	locationController := &LocationController{DB: db}
	router.POST("/api/location/:from/:to", 	locationController.GetLocations)

	server := negroni.New()
	server.Use(negroni.NewLogger())
	server.Use(negroni.NewRecovery())
	server.Use(negroni.HandlerFunc(HeaderMiddleware))
	server.Use(negroni.HandlerFunc(AuthenticationMiddleware))
	server.UseHandler(router)

	log.Println("Starting server on http://" + config.Host + ":" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, server))
}

/*func Index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	http.ServeFile(res, req, config.PublicDir+"/index.html")
}

func NotFoundHandler(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, config.PublicDir+"/404.html")
}*/

func HeaderMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	res.Header().Set("x-frame-options", "SAMEORIGIN")

	if strings.HasPrefix(req.URL.Path, "/api") {

		res.Header().Set("Content-Type", "application/json")

	} else {

		res.Header().Set("x-content-type-options", "nosniff")
		res.Header().Set("x-xss-protection", "1; mode=block")
	}

	next(res, req)
}

func AuthenticationMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	if strings.HasPrefix(req.URL.Path, "/api") {

		// get token from request
		// check if token is valid and find corresponding user
		// append user to request


		if true {

			next(res, req)

		} else {

			log.Println("Unauthorized access attempt by ", req.RemoteAddr, req.UserAgent())
			http.Error(res, "Unauthorized", 401)
		}

	} else {

		next(res, req)
	}
}
