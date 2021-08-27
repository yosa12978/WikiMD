package helpers

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func RespondMessage(w http.ResponseWriter, code int, message string) {
	res := map[string]interface{}{"status_code": code, "message": message}
	RespondJSON(w, code, res)
}
