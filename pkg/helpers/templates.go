package helpers

import (
	"html/template"
	"net/http"

	"github.com/yosa12978/WikiMD/internal/config"
	"github.com/yosa12978/WikiMD/internal/web/midware"
)

func RenderTmpl(w http.ResponseWriter, r *http.Request, path string, p interface{}) {
	temps := []string{
		"./templates/" + path + ".html",
		"./templates/blocks/header.html",
		"./templates/blocks/footer.html",
	}
	t, err := template.ParseFiles(temps...)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	s, err := midware.Store.Get(r, "user_store")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	username := s.Values["username"]
	role := s.Values["role"]
	auth := s.Values["authenticated"]

	if username == nil || role == nil || auth == nil {
		username = ""
		role = ""
		auth = false
	}

	complp := struct {
		Info          config.Wiki
		P             interface{}
		User          string
		Role          string
		Authenticated bool
	}{config.GetConfig().Wiki, p, username.(string), role.(string), auth.(bool)}
	t.ExecuteTemplate(w, path, complp)
}
