package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/moritzschramm/location-tracker-server/model"

	"github.com/julienschmidt/httprouter"
)

// create new device in database
func (controller *Controller) NewDevice(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// check authorization
	if !controller.CheckIfAdmin(res, req) {
		return
	}

	// check password
	password := req.FormValue("password")
	if len(password) < 12 {
		log.Println("Error creating device: password too short")
		http.Error(res, "Password too short", 400)
		return
	}

	// create a new device
	device, err := model.MakeDevice(controller.DB, []byte(password))
	if err != nil {
		log.Println("Error creating device: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	// update mosquitto passwd file
	err = controller.Mqtt.AddUser(device.UUID.String(), password)
	if err != nil {
		log.Println("Error adding user to mosquitto: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	// generate json response
	deviceJson, err := json.Marshal(device)
	if err != nil {
		log.Println("Json: Error creating device: ", err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	res.Write(deviceJson)
}

// delete device from database
func (controller *Controller) DeleteDevice(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	// check authorization
	if !controller.CheckIfAdmin(res, req) {
		return
	}

	// delete device
	uid := params.ByName("uid")
	err := model.DeleteDeviceByUUID(controller.DB, uid)

	if err != nil {

		// delete device from mqtt passwd file
		err := controller.Mqtt.DeleteUser(uid)
		if err != nil {
			log.Println("Error deleting user from mosquitto: ", err.Error())
			http.Error(res, "Internal Server Error", 500)
			return
		}

		res.WriteHeader(http.StatusOK)

	} else {

		log.Println("Error deleting device from database: ", err.Error())
		res.WriteHeader(http.StatusNotFound)
	}
}
