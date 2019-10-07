package api

import (
	"github.com/gorilla/mux"
)

// NewRouter initializes the API routes
func NewRouter(h *handler) *mux.Router {

	if h == nil {
		h = &handler{}
	}

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()
	router.HandleFunc("/stats/downloads/{repo-name}", h.getDownloadsHandler).Methods("GET")
	/*
		// Uncomment this later when we need a variable length option for the returned data
		// Requires modifications in handler to deal with this.
		router.HandleFunc("/stats/downloads/{repo-name}", h.getDownloadsHandler).
			Queries("limit", "{limit}").Methods("GET")
	*/

	return router
}
