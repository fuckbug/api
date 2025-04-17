package handlers

import (
	"fmt"
	"net/http"

	"github.com/fuckbug/api/internal/modules/app"
	"github.com/gorilla/mux"
)

type appHandler struct {
	logger  Logger
	service app.Service
}

func RegisterAppHandlers(
	r *mux.Router,
	logger Logger,
	service app.Service,
) {
	h := &appHandler{
		logger:  logger,
		service: service,
	}

	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)
}

func (h *appHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := h.service.Health(r.Context())

	_, err := w.Write(response)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Health - response error: %s", err))
	}
}
