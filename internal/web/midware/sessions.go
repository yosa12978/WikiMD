package midware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yosa12978/WikiMD/internal/pkg/models"
	"github.com/yosa12978/WikiMD/internal/pkg/repositories"
)

var Store *sessions.CookieStore = sessions.NewCookieStore([]byte("1010101010101010"))

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := Store.Get(r, "user_store")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		username := s.Values["username"]
		if username == nil {
			http.Redirect(w, r, "/auth/login", 301)
			return
		}
		_, err = repositories.NewUserRepository().ReadUser(username.(string))
		if err != nil {
			http.Redirect(w, r, "/auth/login", 301)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ModOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := Store.Get(r, "user_store")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		username := s.Values["username"]
		if username == nil {
			http.Error(w, "403 Forbidden", 403)
			return
		}
		user, err := repositories.NewUserRepository().ReadUser(username.(string))
		if err != nil {
			http.Error(w, "403 Forbidden", 403)
			return
		}
		if user.Role == models.MODER_ROLE {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "403 Forbidden", 403)
	})
}
