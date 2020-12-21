package render

import (
	"net/http"
)

type Render interface {
	Success(w http.ResponseWriter, data interface{}) error
	Error(w http.ResponseWriter, status int, data interface{}) error
}
