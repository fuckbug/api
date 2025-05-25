package handlers

import (
	"net/http"
	"strconv"

	"github.com/fuckbug/api/internal/middleware"
	"github.com/fuckbug/api/internal/modules/project"
	"github.com/fuckbug/api/pkg/httputils"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type projectHandler struct {
	logger   Logger
	validate *v.Validate
	service  project.Service
}

func RegisterProjectHandlers(
	r *mux.Router,
	logger Logger,
	service project.Service,
	jwtKey []byte,
) {
	h := &projectHandler{
		logger:   logger,
		validate: v.New(),
		service:  service,
	}

	routerV1 := r.PathPrefix("/v1/projects").Subrouter()
	routerV1.Use(middleware.Auth(jwtKey))

	routerV1.HandleFunc("", h.Create).Methods(http.MethodPost)
	routerV1.HandleFunc("", h.GetAll).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.GetByID).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}/dsn", h.GetDSNByID).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
	routerV1.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// GetByID godoc
// @Summary Get a project by ID
// @Description Get a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} project.Entity
// @Param id path string true "Project ID"
// @Security BearerAuth
// @Router /v1/projects/{id} [get].
func (h *projectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	entity, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusNotFound, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, entity)
}

// GetDSNByID godoc
// @Summary Get a project DSN
// @Description Get a project DSN
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Param id path string true "Project ID"
// @Security BearerAuth
// @Router /v1/projects/{id}/dsn [get].
func (h *projectHandler) GetDSNByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	dsn, err := h.service.GetDSNByID(r.Context(), id)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusNotFound, err.Error())
		return
	}

	result := struct {
		Dsn string `json:"dsn"`
	}{
		Dsn: dsn,
	}

	httputils.RespondWithJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all projects
// @Description Retrieves a list of all projects from the system
// @Tags projects
// @Accept  json
// @Produce  json
// @Param sort query string false "Sort order (asc or desc)" default(desc) Enums(asc, desc)
// @Param limit query int false "Items per page" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} project.EntityList "Successfully retrieved list of projects"
// @Security BearerAuth
// @Router /v1/projects [get].
func (h *projectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = httputils.DefaultLimit
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = httputils.DefaultOffset
	}

	sortOrder := queryParams.Get("sort")
	if sortOrder != httputils.SortAsc && sortOrder != httputils.SortDesc {
		sortOrder = httputils.DefaultSort
	}

	params := project.GetAllParams{
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}

	projects, totalCount, err := h.service.GetAll(r.Context(), params)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, httputils.NewListResponse(totalCount, projects))
}

// Create godoc
// @Summary Create a new project entry
// @Description Creates a new project entry in the system
// @Tags projects
// @Accept  json
// @Produce json
// @Param   request body project.Create true "Project entry creation data"
// @Success 201 {object} project.Entity "Successfully created project entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /v1/projects [post].
func (h *projectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req project.Create

	if err := httputils.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputils.HandleValidatorError(w, err)
		return
	}

	entity, err := h.service.Create(r.Context(), &req)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusCreated, entity)
}

// Update godoc
// @Summary Update a project entry
// @Description Updates an existing project entry
// @Tags projects
// @Accept  json
// @Produce json
// @Param   id path string true "Project entry ID"
// @Param   request body project.Update true "Project update data"
// @Success 200 {object} project.Entity "Successfully updated project entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Project entry not found"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /v1/projects/{id} [put].
func (h *projectHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req project.Update
	if err := httputils.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputils.HandleValidatorError(w, err)
		return
	}

	entity, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, entity)
}

// Delete godoc
// @Summary Delete a project entry
// @Description Delete a project entry by its ID
// @Tags projects
// @Accept  json
// @Produce  json
// @Param id path string true "Project entry ID"
// @Success 204 "No Content"
// @Failure 400 {object} string "Bad Request - when ID is not provided"
// @Failure 500 {object} string "Internal Server Error - when something goes wrong"
// @Security BearerAuth
// @Router /v1/projects/{id} [delete]
func (h *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
