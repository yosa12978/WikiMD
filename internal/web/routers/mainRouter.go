package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/web/handlers"
)

func InitMainRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	mainHandler := handlers.NewMainHandler()

	api.HandleFunc("/", mainHandler.GetWikiInfo).Methods("GET")
	return router
}
