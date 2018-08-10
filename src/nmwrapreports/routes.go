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
		"Login",
		"POST",
		"/login",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		Logout,
	},
	Route{
		"ResetPassword",
		"POST",
		"/resetpassword",
		ResetPassword,
	},
	Route{
		"ChangePassword",
		"POST",
		"/changepassword",
		ChangePassword,
	},
	Route{
		"CheckReset",
		"POST",
		"/checkreset",
		CheckReset,
	},
	Route{
		"GetReportFromUpload",
		"POST",
		"/reportupload",
		GetReportFromUpload,
	},
	Route{
		"LoggedIn",
		"GET",
		"/loggedin",
		LoggedIn,
	},
	Route{
		"History",
		"GET",
		"/history",
		History,
	},
	Route{
		"DeleteHistory",
		"POST",
		"/deletehistory",
		DeleteHistory,
	},
	Route{
		"CreateUser",
		"POST",
		"/createuser",
		CreateUser,
	},
	Route{
		"GetReport",
		"GET",
		"/getreport/{key}/{fname}",
		GetReport,
	},
}
