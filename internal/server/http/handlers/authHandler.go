package handlers

import (
	"net/http"

	"github.com/fuckbug/api/internal/modules/users"
	"github.com/fuckbug/api/pkg/httputils"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type authHandler struct {
	logger   Logger
	validate *v.Validate
	service  users.Service
}

func RegisterAuthHandlers(
	r *mux.Router,
	logger Logger,
	service users.Service,
) {
	h := &authHandler{
		logger:   logger,
		validate: v.New(),
		service:  service,
	}

	routerV1 := r.PathPrefix("/v1").Subrouter()
	routerV1.HandleFunc("/signup", h.Signup).Methods(http.MethodPost)
	routerV1.HandleFunc("/login", h.Login).Methods(http.MethodPost)
}

// Signup godoc
// @Summary Signup
// @Description Signup
// @Tags auth
// @Accept  json
// @Produce json
// @Param   request body users.Signup true "Signup"
// @Success 201 {object} int "Successfully signup"
// @Failure 400 {object} string "Invalid input data"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/signup [post].
func (h *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req users.Signup
	if err := httputils.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputils.HandleValidatorError(w, err)
		return
	}

	err := h.service.Signup(r.Context(), &req)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusCreated, 1)
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce json
// @Param   request body users.Login true "Login"
// @Success 201 {object} int "Successfully login"
// @Failure 400 {object} string "Invalid input data"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/login [post].
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req users.Login
	if err := httputils.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputils.HandleValidatorError(w, err)
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusCreated, res)
}
