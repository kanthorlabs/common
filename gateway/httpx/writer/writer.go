package writer

import (
	"net/http"

	"github.com/kanthorlabs/common/utils"
)

// E is shortcut of error data type
type E struct {
	Error string `json:"error" example:"KANTHOR.SYSTEM.ERROR"`
}

// M is shortcut of generic data type
type M map[string]any

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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

func ErrConflict(w http.ResponseWriter, data any) {
	Write(w, http.StatusConflict, data)
}

func ErrUnauthorized(w http.ResponseWriter, data any) {
	Write(w, http.StatusUnauthorized, data)
}

func ErrNotFound(w http.ResponseWriter, data any) {
	Write(w, http.StatusNotFound, data)
}

func ErrUnknown(w http.ResponseWriter, data any) {
	Write(w, http.StatusInternalServerError, data)
}
