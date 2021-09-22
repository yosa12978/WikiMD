package handlers

import (
	"net/http"

	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type IErrorHandler interface {
	Error404Handler(w http.ResponseWriter, r *http.Request)
}

type ErrorHandler struct{}

func NewErrorHandler() IErrorHandler {
	return &ErrorHandler{}
}

func (eh *ErrorHandler) Error404Handler(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTmpl(w, r, "404", nil)
}
