package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Route struct for basic structure of a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes is just a slice
type Routes []Route

// NewRouter function is for setting up basic routes
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Version",
		"GET",
		"/version/",
		Version,
	},
	Route{
		"GetPOSTGeom",
		"GET",
		"/postgeom",
		GetPOSTGeom,
	},
	Route{
		"POSTGeom",
		"POST",
		"/postgeom",
		POSTGeom,
	},
	Route{
		"GetReport",
		"GET",
		"/getreport/{key}/{fname}",
		GetReport,
	},
}
