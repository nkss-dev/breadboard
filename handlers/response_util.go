package handlers

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	data := map[string]interface{}{"data": payload}
	response, err := json.Marshal(data)
	if err != nil {
		RespondError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(response)
}

func RespondError(w http.ResponseWriter, status int, detail string) {
	response, _ := json.Marshal(map[string][]Error{
		"errors": {{Title: http.StatusText(status), Detail: detail}},
	})
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(response)
}
