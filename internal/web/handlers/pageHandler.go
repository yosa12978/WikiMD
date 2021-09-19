package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/pkg/dto"
	"github.com/yosa12978/WikiMD/internal/pkg/repositories"
	"github.com/yosa12978/WikiMD/internal/web/midware"
	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type IPageHandler interface {
	GetPage(w http.ResponseWriter, r *http.Request)
	GetPages(w http.ResponseWriter, r *http.Request)
	CreatePageGet(w http.ResponseWriter, r *http.Request)
	CreatePagePost(w http.ResponseWriter, r *http.Request)
	UpdatePageGet(w http.ResponseWriter, r *http.Request)
	UpdatePagePost(w http.ResponseWriter, r *http.Request)
	DeletePage(w http.ResponseWriter, r *http.Request)
	SearchPage(w http.ResponseWriter, r *http.Request)
}

type PageHandler struct{}

func NewPageHandler() IPageHandler {
	return &PageHandler{}
}

func (ph *PageHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	vrs := mux.Vars(r)
	page, err := repositories.NewPageRepository().ReadPage(vrs["id"])
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	helpers.RenderTmpl(w, r, "page", page)
}

func (ph *PageHandler) GetPages(w http.ResponseWriter, r *http.Request) {
	pages := repositories.NewPageRepository().GetPages()
	helpers.RenderTmpl(w, r, "pages", pages)
}

func (ph *PageHandler) CreatePageGet(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTmpl(w, r, "createpage", nil)
}

func (ph *PageHandler) CreatePagePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s, _ := midware.Store.Get(r, "user_store")
	cpd := dto.CreatePageDTO{Name: r.FormValue("name"), Body: r.FormValue("body")}
	err := repositories.NewPageRepository().CreatePage(&cpd, s.Values["username"].(string))
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	http.Redirect(w, r, "/pages/", 301)
}

func (ph *PageHandler) UpdatePageGet(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	commit, err := repositories.NewPageRepository().ReadPage(v["id"])
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	helpers.RenderTmpl(w, r, "updpage", commit)
}

func (ph *PageHandler) UpdatePagePost(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	r.ParseForm()
	s, err := midware.Store.Get(r, "user_store")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	username := s.Values["username"].(string)
	cdto := dto.CreateCommitDTO{
		Name:   r.FormValue("name"),
		Body:   r.FormValue("body"),
		PageID: v["id"],
		User:   username,
	}
	err = repositories.NewCommitRepository().CreateCommit(cdto, username)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	http.Redirect(w, r, "/page/id/"+v["id"], 301)
}

func (ph *PageHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, err := midware.Store.Get(r, "user_store")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	username := s.Values["username"]
	if username == nil {
		username = ""
	}
	err = repositories.NewPageRepository().DeletePage(v["id"], username.(string))
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	http.Redirect(w, r, "/pages", 301)
}

func (ph *PageHandler) SearchPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		helpers.RenderTmpl(w, r, "search", nil)
		return
	}
	pages := repositories.NewPageRepository().SearchPages(query)
	helpers.RenderTmpl(w, r, "search", pages)
}
