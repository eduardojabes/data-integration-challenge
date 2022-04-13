package company

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range hook {
		var handler http.Handler
		handler = route.HandlerFunc

		var muxRoute *mux.Route
		muxRoute = router.Methods(route.Method)
		muxRoute.
			Name(route.Name).
			Handler(handler)
	}

	return router
}
