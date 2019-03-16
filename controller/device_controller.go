package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

type DeviceController struct {
	DB *sql.DB
}

func (controller *DeviceController) NewDevice(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// create a new device
	password := req.FormValue("password")
	device := model.MakeDevice(controller.DB, password)

	// update mosquitto passwd and acl file

	deviceJson, err := json.Marshal(device)
	if err != nil {
		log.Println("Error creating note: ", err)
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(deviceJson)
}

func (controller *DeviceController) DeleteDevice(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	uid := params.ByName("uid")

	err := model.DeleteDeviceByUUID(controller.DB, uid)

	if err != nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
