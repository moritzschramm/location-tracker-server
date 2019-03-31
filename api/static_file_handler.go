package api

import (
	"net/http"
	"strings"

	"github.com/moritzschramm/location-tracker-server/config"

	"github.com/julienschmidt/httprouter"
)

type StaticFileHandler struct {
	Config config.Config
}

func (h *StaticFileHandler) SetupRoutes(router *httprouter.Router) {

	router.NotFound = http.HandlerFunc(h.NotFoundHandler)

	router.GET("/", h.ServeFile("index.html"))
	router.GET("/favicon.ico", h.ServeFile("favicon.ico"))
	router.GET("/robots.txt", h.ServeFile("robots.txt"))

	router.ServeFiles("/assets/*filepath", http.Dir(h.Config.PublicDir+"/dist/assets"))
}

func (h *StaticFileHandler) NotFoundHandler(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, h.Config.PublicDir+"/404.html")
}

func (h *StaticFileHandler) ServeFile(filename string) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

		http.ServeFile(res, req, h.Config.PublicDir+"/dist/"+filename)
	}
}

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
