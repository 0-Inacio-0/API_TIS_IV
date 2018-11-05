package gyms

import (
	"github.com/0-Inacio-0/API_TIS_IV/logger"
	"github.com/gorilla/mux"
	"net/http"
)

var controller = &Controller{}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes of the API
type Routes []Route

var routes = Routes{
	Route{
		"Home",
		"GET",
		"/",
		controller.Home,
	},
	Route{
		"GetGyms",
		"GET",
		"/gyms",
		controller.GetGyms,
	},
	Route{
		"PostScore",
		"post",
		"/score",
		controller.PostScore,
	},
}

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
