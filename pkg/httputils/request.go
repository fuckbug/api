package httputils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	DefaultLimit  = 50
	DefaultOffset = 0
	DefaultSort   = SortDesc
	SortAsc       = "asc"
	SortDesc      = "desc"
)

func DecodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if r.Body == nil {
		RespondWithError(w, http.StatusBadRequest, "Request body is required", nil)
		return errors.New("empty request body")
	}
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		RespondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json", nil)
		return errors.New("invalid content type")
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return err
	}

	return nil
}
