package handlers

import (
	"net/http"

	"github.com/yosa12978/WikiMD/internal/pkg/repositories"
	"github.com/yosa12978/WikiMD/internal/web/midware"
	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type IAuthHandler interface {
	LoginUserPost(w http.ResponseWriter, r *http.Request)
	CreateUserPost(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
	LoginUserGet(w http.ResponseWriter, r *http.Request)
	CreateUserGet(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
}

func NewAuthHandler() IAuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) LoginUserGet(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTmpl(w, r, "login", nil)
}

func (ah *AuthHandler) LoginUserPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	usd, err := repositories.NewUserRepository().LogInUser(username, password)
	if err != nil {
		w.Write([]byte("user does not exist"))
		return
	}
	s, err := midware.Store.Get(r, "user_store")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	s.Values["username"] = usd.Username
	s.Values["role"] = string(usd.Role)
	s.Values["authenticated"] = true
	err = s.Save(r, w)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", 301)
}

func (ah *AuthHandler) CreateUserGet(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTmpl(w, r, "signup", nil)
}

func (ah *AuthHandler) CreateUserPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	err := repositories.NewUserRepository().CreateUser(username, password, email)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	http.Redirect(w, r, "/auth/login", 301)
}

func (ah *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	s, err := midware.Store.Get(r, "user_store")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	delete(s.Values, "username")
	delete(s.Values, "role")
	delete(s.Values, "authenticated")
	err = s.Save(r, w)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/auth/login", 301)
}
