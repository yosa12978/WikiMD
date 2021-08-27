package routers

import (
	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/web/handlers"
)

func InitMainRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	mainHandler := handlers.NewMainHandler()

	api.HandleFunc("/", mainHandler.GetWikiInfo).Methods("GET")
	return router
}
