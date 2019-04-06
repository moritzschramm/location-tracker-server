package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// setup all static routes for httprouter
// index, favicon and robots.txt are set explicitly
// all files in /assets/* get served from <PublicDir>/dist/assets
func (c *Controller) SetupStaticRoutes(router *httprouter.Router) {

	router.NotFound = http.HandlerFunc(c.NotFoundHandler)

	router.GET("/", c.ServeFile("index.html"))
	router.GET("/favicon.ico", c.ServeFile("favicon.ico"))
	router.GET("/robots.txt", c.ServeFile("robots.txt"))

	router.ServeFiles("/assets/*filepath", http.Dir(c.Config.PublicDir+"/dist/assets"))
}

func (c *Controller) NotFoundHandler(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, c.Config.PublicDir+"/404.html")
}

// helper function to serve files
func (c *Controller) ServeFile(filename string) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

		http.ServeFile(res, req, c.Config.PublicDir+"/dist/"+filename)
	}
}

// header middleware sets correct headers depending on request URI
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
