package render

import (
	"encoding/json"
	"net/http"
)

type JsonRender struct{}

func NewJsonRender() *JsonRender {
	return &JsonRender{}
}

func (r *JsonRender) Success(w http.ResponseWriter, data interface{}) error {
	return r.encode(w, http.StatusOK, data)
}

func (r *JsonRender) Error(w http.ResponseWriter, status int, data interface{}) error {
	return r.encode(w, status, data)
}

func (r *JsonRender) encode(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
