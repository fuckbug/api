package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getProjectIDAndKey(r *http.Request) (projectID string, err error) {
	vars := mux.Vars(r)
	projectID = vars["projectID"]
	key := vars["key"]
	if projectID == "" || key == "" {
		return "", fmt.Errorf("invalid ingest")
	}
	return projectID, nil
}
