package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

// return battery info for specified time frame
func (controller *Controller) GetBatteryInfo(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	device := req.Context().Value("device").(*model.Device)

	// parse from and to datetimes
	from, err := time.Parse(time.RFC3339, params.ByName("from"))
	if err != nil {
		http.Error(res, "'from' parameter malformed", 400)
	}

	to, err := time.Parse(time.RFC3339, params.ByName("to"))
	if err != nil {
		http.Error(res, "'to' parameter malformed", 400)
	}

	batteryInfo, err := model.GetBatteryInfo(controller.DB, device.DeviceId, from, to)
	if err != nil {
		log.Println("Error retrieving battery info: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	batteryInfoJson, err := json.Marshal(batteryInfo)
	if err != nil {
		log.Println("Json: Error parsing battery info: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(batteryInfoJson)
}
