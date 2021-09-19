package midware

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func CountTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcname := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
		t := time.Now()
		next.ServeHTTP(w, r)
		str := fmt.Sprintf("%s - %s", funcname, time.Since(t))
		log.Println(str)
	})
}
