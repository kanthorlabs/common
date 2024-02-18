package writer

import (
	"net/http"

	"github.com/kanthorlabs/common/utils"
)

// M is shortcu of generic data type
type M map[string]any

func Write(w http.ResponseWriter, status, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(utils.Stringify(data)))
}

func Ok(w http.ResponseWriter, data any) {
	Write(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	Write(w, http.StatusCreated, data)
}

func ErrBadRequest(w http.ResponseWriter, data any) {
	Write(w, http.StatusBadRequest, data)
}

func ErrUnauthorized(w http.ResponseWriter, data any) {
	Write(w, http.StatusUnauthorized, data)
}

func ErrNotFound(w http.ResponseWriter, data any) {
	Write(w, http.StatusNotFound, data)
}

func Err(w http.ResponseWriter, data any) {
	Write(w, http.StatusInternalServerError, data)
}
