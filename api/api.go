package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	handler := newHandler()

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()
	router.HandleFunc("/stats/downloads/{repo-name}", handler.getDownloadsHandler).Methods("GET")
	router.HandleFunc("/stats/downloads/{repo-name}", handler.getDownloadsHandler).
		Queries("limit", "{limit}").Methods("GET")

	return router
}
