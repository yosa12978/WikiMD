package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/web/handlers"
	"github.com/yosa12978/WikiMD/internal/web/midware"
)

func InitMainRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(midware.CountTime)

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// mainHandler := handlers.NewMainHandler()
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	pageHandler := handlers.NewPageHandler()
	commitHandler := handlers.NewCommitHandler()
	errorHandler := handlers.NewErrorHandler()

	router.HandleFunc("/page/id/{id}", pageHandler.GetPage).Methods("GET")
	router.HandleFunc("/", pageHandler.GetPages).Methods("GET")
	router.Handle("/page/create", midware.LoginRequired(http.HandlerFunc(pageHandler.CreatePageGet))).Methods("GET")
	router.Handle("/page/create", midware.LoginRequired(http.HandlerFunc(pageHandler.CreatePagePost))).Methods("POST")
	router.Handle("/page/update/{id}", midware.LoginRequired(http.HandlerFunc(pageHandler.UpdatePageGet))).Methods("GET")
	router.Handle("/page/update/{id}", midware.LoginRequired(http.HandlerFunc(pageHandler.UpdatePagePost))).Methods("POST")
	router.HandleFunc("/pages/search", pageHandler.SearchPage).Methods("GET")
	router.HandleFunc("/commits/{page_id}", commitHandler.GetPageCommits).Methods("GET")
	router.HandleFunc("/commit/{id}", commitHandler.GetCommit).Methods("GET")
	router.Handle("/commit/delete/{id}", midware.ModOnly(http.HandlerFunc(commitHandler.DeleteCommit))).Methods("GET")
	router.Handle("/page/delete/{id}", midware.ModOnly(http.HandlerFunc(pageHandler.DeletePage))).Methods("GET")

	router.HandleFunc("/user/{username}", userHandler.GetUser).Methods("GET")

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authHandler.LoginUserGet).Methods("GET")
	auth.HandleFunc("/login", authHandler.LoginUserPost).Methods("POST")
	auth.HandleFunc("/signup", authHandler.CreateUserGet).Methods("GET")
	auth.HandleFunc("/signup", authHandler.CreateUserPost).Methods("POST")
	auth.HandleFunc("/logout", authHandler.LogoutUser).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(errorHandler.Error404Handler)

	return router
}
