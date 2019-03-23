package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/moritzschramm/location-tracker-server/config"
	"github.com/moritzschramm/location-tracker-server/mqtt"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func SetupAPI(db *sql.DB, mqttClient MQTT.Client, config config.Config) {

	// setup router
	router := httprouter.New()

	// serve static files
	staticFileHandler := &StaticFileHandler{config}
	router.NotFound = http.HandlerFunc(staticFileHandler.NotFoundHandler)
	router.GET("/", staticFileHandler.ServeSinglePageApplication)
	router.ServeFiles("/"+config.PublicDir+"/assets/*filepath", http.Dir(config.PublicDir+"/assets"))

	// api routes
	deviceController := &DeviceController{DB: db, Mqtt: mqtt.User{config.MQTT}}
	router.POST("/api/device/new", deviceController.NewDevice)
	router.POST("/api/device/delete/:uid", deviceController.DeleteDevice)

	locationController := &LocationController{DB: db}
	router.POST("/api/location/:from/:to", locationController.GetLocations)

	authController := &AuthController{DB: db}
	router.POST("/api/login", authController.Login)
	router.POST("/api/logout", authController.Logout)
	router.POST("/api/refresh", authController.TokenRefresh)

	// setup negroni middleware
	server := negroni.New()
	server.Use(negroni.NewLogger())
	server.Use(negroni.NewRecovery())
	server.Use(negroni.HandlerFunc(HeaderMiddleware))
	server.Use(negroni.HandlerFunc(authController.AuthenticationMiddleware))
	server.UseHandler(router)

	log.Println("Starting server on http://" + config.Host + ":" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, server))
}
