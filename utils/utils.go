package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Result struct {
	SearchResults []string `json:"results"`
	Correction    string   `json:"correction"`
}

// get the app port from env variable ...
func GetAppPort() string {
	appPort := os.Getenv("PORT")
	if appPort == "" {
		return "3001"
	}
	return appPort
}

// fatal error check ...
func ErrorFatalCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// json encoding for search result ...
func EncodeResult(result []string, correction string) (*bytes.Buffer, error) {
	resp := Result{
		SearchResults: result,
		Correction:    correction,
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err := enc.Encode(resp)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// write back result to response writer ...
func WriteResponse(statusCode int, result []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}
