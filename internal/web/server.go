package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/web/routers"
)

func Run(port int) {
	router := routers.InitMainRouter()
	go func(prt int, routr *mux.Router) {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", prt), routr)
	}(port, router)
}
