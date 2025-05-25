package handlers

import (
	"net/http"
	"strconv"

	"github.com/fuckbug/api/internal/middleware"
	"github.com/fuckbug/api/internal/modules/errors"
	"github.com/fuckbug/api/pkg/httputils"
	"github.com/fuckbug/api/pkg/utils"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type errorHandler struct {
	logger   Logger
	validate *v.Validate
	service  errors.Service
}

func RegisterErrorHandlers( //nolint:dupl
	r *mux.Router,
	logger Logger,
	service errors.Service,
	jwtKey []byte,
) {
	h := &errorHandler{
		logger:   logger,
		validate: v.New(),
		service:  service,
	}

	r.HandleFunc("/ingest/{projectID}:{key}/errors", h.Create).Methods(http.MethodPost)

	routerV1 := r.PathPrefix("/v1/errors").Subrouter()
	routerV1.Use(middleware.Auth(jwtKey))

	routerV1.HandleFunc("", h.GetAll).Methods(http.MethodGet)
	routerV1.HandleFunc("/stats", h.GetStats).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.GetByID).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
	routerV1.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// GetByID godoc
// @Summary Get an error by ID
// @Description Get an error by ID
// @Tags errors
// @Accept json
// @Produce json
// @Success 200 {object} errors.Entity
// @Param id path string true "Error ID"
// @Security BearerAuth
// @Router /v1/errors/{id} [get].
func (h *errorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

// GetAll godoc
// @Summary Get all errors
// @Description Retrieves a list of all errors from the system
// @Tags errors
// @Accept json
// @Produce json
// @Param projectId query string false "Project ID"
// @Param groupId query string false "Group ID"
// @Param timeFrom query int false "Time errors from"
// @Param timeTo query int false "Time errors to"
// @Param search query string false "Search in message field"
// @Param sort query string false "Sort order (asc or desc)" default(desc) Enums(asc, desc)
// @Param limit query int false "Items per page" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} errors.EntityList "Successfully retrieved list of errors"
// @Security BearerAuth
// @Router /v1/errors [get].
func (h *errorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = httputils.DefaultLimit
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = httputils.DefaultOffset
	}

	projectID := queryParams.Get("projectId")
	groupID := queryParams.Get("groupId")

	timeFrom, err := utils.ParseTimeParam(queryParams.Get("timeFrom"))
	if err != nil {
		timeFrom = 0
	}

	timeTo, err := utils.ParseTimeParam(queryParams.Get("timeTo"))
	if err != nil {
		timeTo = 0
	}

	sortOrder := queryParams.Get("sort")
	if sortOrder != httputils.SortAsc && sortOrder != httputils.SortDesc {
		sortOrder = httputils.DefaultSort
	}

	search := queryParams.Get("search")

	params := errors.GetAllParams{
		FilterParams: errors.FilterParams{
			ProjectID:   projectID,
			Fingerprint: groupID,
			TimeFrom:    utils.SecondsToMilliseconds(timeFrom),
			TimeTo:      utils.SecondsToMilliseconds(timeTo),
			Search:      search,
		},
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}

	entities, totalCount, err := h.service.GetAll(r.Context(), params)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, httputils.NewListResponse(totalCount, entities))
}

// GetStats godoc
// @Summary Get errors stats
// @Description Retrieves a stats of all errors from the system
// @Tags errors
// @Accept json
// @Produce json
// @Param projectId query string true "Project ID"
// @Param groupId query string false "Group ID"
// @Success 200 {object} errors.Stats "Successfully retrieved stats of errors"
// @Security BearerAuth
// @Router /v1/errors/stats [get].
func (h *errorHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	projectID := queryParams.Get("projectId")
	groupID := queryParams.Get("groupId")

	stats, err := h.service.GetStats(r.Context(), projectID, groupID)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, stats)
}

// Create godoc
// @Summary Create a new error entry
// @Description Creates a new error entry in the system
// @Tags ingest
// @Accept  json
// @Produce json
// @Param        projectID   path      string  true  "Project ID"
// @Param        key         path      string  true  "Public key"
// @Param   request body errors.Create true "Error entry creation data"
// @Success 201 {object} errors.Entity "Successfully created error entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 500 {object} string "Internal server error"
// @Router /ingest/{projectID}:{key}/errors [post].
func (h *errorHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := getProjectIDAndKey(r)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req errors.Create
	if err := httputils.DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httputils.HandleValidatorError(w, err)
		return
	}

	req.ProjectID = projectID

	entity, err := h.service.Create(r.Context(), &req)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusCreated, entity)
}

// Update godoc
// @Summary Update an error entry
// @Description Updates an existing error entry
// @Tags errors
// @Accept  json
// @Produce json
// @Param   id path string true "Error entry ID"
// @Param   request body errors.Update true "Error update data"
// @Success 200 {object} errors.Entity "Successfully updated error entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Error entry not found"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /v1/errors/{id} [put].
func (h *errorHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req errors.Update
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
// @Summary Delete an error entry
// @Description Delete an error entry by its ID
// @Tags errors
// @Accept  json
// @Produce  json
// @Param id path string true "Error entry ID"
// @Success 204 "No Content"
// @Failure 400 {object} string "Bad Request - when ID is not provided"
// @Failure 500 {object} string "Internal Server Error - when something goes wrong"
// @Security BearerAuth
// @Router /v1/errors/{id} [delete]
func (h *errorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
