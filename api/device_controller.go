package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/moritzschramm/location-tracker-server/model"
	"github.com/moritzschramm/location-tracker-server/mqtt"

	"github.com/julienschmidt/httprouter"
)

type DeviceController struct {
	DB *sql.DB
	Mqtt mqtt.User
}

func (controller *DeviceController) NewDevice(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// create a new device
	password := req.FormValue("password")
	device, err := model.MakeDevice(controller.DB, []byte(password))
	if err != nil {
		log.Println("Error creating device: ", err)
		http.Error(res, "Internal Server Error", 500)
		return
	}

	// update mosquitto passwd file
	err = controller.Mqtt.AddUser(device.UUID.String(), password)
	if err != nil {
		log.Println("Error adding user to mosquitto: ", err)
		http.Error(res, "Internal Server Error", 500)
		return
	}

	// generate json response
	deviceJson, err := json.Marshal(device)
	if err != nil {
		log.Println("Json: Error creating device: ", err)
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(deviceJson)
}

func (controller *DeviceController) DeleteDevice(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	uid := params.ByName("uid")

	err := model.DeleteDeviceByUUID(controller.DB, uid)

	if err != nil {

		// delete device from mqtt passwd file
		err := controller.Mqtt.DeleteUser(uid)
		if err != nil {
			log.Println("Error deleting user from mosquitto: ", err)
		}

		res.WriteHeader(http.StatusOK)

	} else {

		res.WriteHeader(http.StatusNotFound)
	}
}
