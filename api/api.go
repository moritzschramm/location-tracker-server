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

// controller handles API calls and serves static files
// has database and configuration access
type Controller struct {
	DB     *sql.DB
	Config config.Config
	Mqtt   mqtt.Config
}

// setup router and middleware and start server
// router handles static files and API calls
func Setup(db *sql.DB, mqttClient MQTT.Client, config config.Config) {

	router := httprouter.New()
	controller := &Controller{
		DB:     db,
		Config: config,
		Mqtt:   config.MQTT,
	}

	// serve static files
	controller.SetupStaticRoutes(router)

	// api routes
	// device
	router.POST("/api/device/new", controller.NewDevice)
	router.POST("/api/device/delete/:uid", controller.DeleteDevice)

	// location
	router.POST("/api/location/:from/:to", controller.GetLocations)

	// authentication
	router.POST("/api/login", controller.Login)
	router.POST("/api/logout", controller.Logout)
	router.POST("/api/refresh", controller.TokenRefresh)

	// setup negroni middleware
	server := negroni.New()
	server.Use(negroni.NewLogger())
	server.Use(negroni.NewRecovery())
	server.Use(negroni.HandlerFunc(HeaderMiddleware))
	server.Use(negroni.HandlerFunc(controller.AuthenticationMiddleware))
	server.UseHandler(router)

	// start TLS encrypted server
	log.Println("Starting server on https://" + config.Host + ":" + config.Port)
	log.Fatal(http.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile, server))
}
