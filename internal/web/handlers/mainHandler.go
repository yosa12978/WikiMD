package handlers

import (
	"net/http"

	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type IMainHandler interface {
	GetWikiInfo(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
}

func NewMainHandler() IMainHandler {
	return &MainHandler{}
}

func (mh *MainHandler) GetWikiInfo(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTmpl(w, r, "index", nil)
}
