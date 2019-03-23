package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

type LocationController struct {
	DB *sql.DB
}

func (controller *LocationController) GetLocations(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	device := req.Context().Value("device").(*model.Device)

	from, err := time.Parse(time.RFC3339, params.ByName("from"))
	if err != nil {
		http.Error(res, "'from' parameter malformed", 400)
	}

	to, err := time.Parse(time.RFC3339, params.ByName("to"))
	if err != nil {
		http.Error(res, "'to' parameter malformed", 400)
	}

	locations, err := model.GetLocations(controller.DB, device.DeviceId, from, to)
	if err != nil {
		log.Println("Error retrieving locations: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	locationsJson, err := json.Marshal(locations)
	if err != nil {
		log.Println("Json: Error retrieving locations: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(locationsJson)
}
