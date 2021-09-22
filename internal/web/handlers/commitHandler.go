package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/WikiMD/internal/pkg/repositories"
	"github.com/yosa12978/WikiMD/pkg/helpers"
)

type ICommitHandler interface {
	GetPageCommits(w http.ResponseWriter, r *http.Request)
	GetCommit(w http.ResponseWriter, r *http.Request)
	DeleteCommit(w http.ResponseWriter, r *http.Request)
}

type CommitHandler struct{}

func NewCommitHandler() ICommitHandler {
	return &CommitHandler{}
}

func (ch *CommitHandler) GetPageCommits(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	commits, err := repositories.NewCommitRepository().GetCommitsByPageID(v["page_id"])
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	helpers.RenderTmpl(w, r, "commits", commits)
}

func (ch *CommitHandler) GetCommit(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	commit, err := repositories.NewCommitRepository().GetCommitByID(v["id"])
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	helpers.RenderTmpl(w, r, "commit", commit)
}

func (ch *CommitHandler) DeleteCommit(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	err := repositories.NewCommitRepository().DeleteCommit(v["id"])
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	http.Redirect(w, r, "/", 301)
}
