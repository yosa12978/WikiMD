package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/pkg/repositories"
	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type IUserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct{}

func NewUserHandler() IUserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ur := repositories.NewUserRepository()
	user, err := ur.ReadUser(mux.Vars(r)["username"])
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	helpers.RenderTmpl(w, r, "user", user)
}
