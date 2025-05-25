package handlers

import (
	"net/http"
	"strconv"

	"github.com/fuckbug/api/internal/middleware"
	"github.com/fuckbug/api/internal/modules/log"
	"github.com/fuckbug/api/pkg/httputils"
	"github.com/fuckbug/api/pkg/utils"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type logHandler struct {
	logger   Logger
	validate *v.Validate
	service  log.Service
}

func RegisterLogHandlers( //nolint:dupl
	r *mux.Router,
	logger Logger,
	service log.Service,
	jwtKey []byte,
) {
	h := &logHandler{
		logger:   logger,
		validate: v.New(),
		service:  service,
	}

	r.HandleFunc("/ingest/{projectID}:{key}/logs", h.Create).Methods(http.MethodPost)

	routerV1 := r.PathPrefix("/v1/logs").Subrouter()
	routerV1.Use(middleware.Auth(jwtKey))

	routerV1.HandleFunc("", h.GetAll).Methods(http.MethodGet)
	routerV1.HandleFunc("/stats", h.GetStats).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.GetByID).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
	routerV1.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// GetByID godoc
// @Summary Get a log by ID
// @Description Get a log by ID
// @Tags logs
// @Accept json
// @Produce json
// @Success 200 {object} log.Entity
// @Param id path string true "Log ID"
// @Security BearerAuth
// @Router /v1/logs/{id} [get].
func (h *logHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
// @Summary Get all logs
// @Description Retrieves a list of all logs from the system
// @Tags logs
// @Accept  json
// @Produce  json
// @Param projectId query string false "Project ID"
// @Param groupId query string false "Group ID"
// @Param timeFrom query int false "Time logs from"
// @Param timeTo query int false "Time logs to"
// @Param level query string false "Filter by log level" Enums(DEBUG, INFO, WARN, ERROR)
// @Param search query string false "Search in message field"
// @Param sort query string false "Sort order (asc or desc)" default(desc) Enums(asc, desc)
// @Param limit query int false "Items per page" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} log.EntityList "Successfully retrieved list of logs"
// @Security BearerAuth
// @Router /v1/logs [get].
func (h *logHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	level := queryParams.Get("level")
	search := queryParams.Get("search")

	params := log.GetAllParams{
		FilterParams: log.FilterParams{
			ProjectID:   projectID,
			Fingerprint: groupID,
			TimeFrom:    timeFrom,
			TimeTo:      timeTo,
			Level:       level,
			Search:      search,
		},
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}

	logs, totalCount, err := h.service.GetAll(r.Context(), params)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, httputils.NewListResponse(totalCount, logs))
}

// GetStats godoc
// @Summary Get logs stats
// @Description Retrieves a stats of all logs from the system
// @Tags logs
// @Accept  json
// @Produce  json
// @Param projectId query string true "Project ID"
// @Param groupId query string false "Group ID"
// @Success 200 {object} log.Stats "Successfully retrieved stats of logs"
// @Security BearerAuth
// @Router /v1/logs/stats [get].
func (h *logHandler) GetStats(w http.ResponseWriter, r *http.Request) {
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
// @Summary Create a new log entry
// @Description Creates a new log entry in the system
// @Tags ingest
// @Accept  json
// @Produce json
// @Param        projectID   path      string  true  "Project ID"
// @Param        key         path      string  true  "Public key"
// @Param   request body log.Create true "Log entry creation data"
// @Success 201 {object} log.Entity "Successfully created log entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 500 {object} string "Internal server error"
// @Router /ingest/{projectID}:{key}/logs [post].
func (h *logHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := getProjectIDAndKey(r)
	if err != nil {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req log.Create
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
// @Summary Update a log entry
// @Description Updates an existing log entry
// @Tags logs
// @Accept  json
// @Produce json
// @Param   id path string true "Log entry ID"
// @Param   request body log.Update true "Log update data"
// @Success 200 {object} log.Entity "Successfully updated log entry"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Log entry not found"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /v1/logs/{id} [put].
func (h *logHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		httputils.RespondWithPlainError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req log.Update
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
// @Summary Delete a log entry
// @Description Delete a log entry by its ID
// @Tags logs
// @Accept  json
// @Produce  json
// @Param id path string true "Log entry ID"
// @Success 204 "No Content"
// @Failure 400 {object} string "Bad Request - when ID is not provided"
// @Failure 500 {object} string "Internal Server Error - when something goes wrong"
// @Security BearerAuth
// @Router /v1/logs/{id} [delete]
func (h *logHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
