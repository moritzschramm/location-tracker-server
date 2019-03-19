package api

import (
	"net/http"

	"github.com/moritzschramm/location-tracker-server/config"

	"github.com/julienschmidt/httprouter"
)

type StaticFileHandler struct {
	Config config.Config
}

func (h *StaticFileHandler) NotFoundHandler(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, h.Config.PublicDir+"/404.html")
}

func (h *StaticFileHandler) ServeSinglePageApplication(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	http.ServeFile(res, req, h.Config.PublicDir+"/index.html")
}