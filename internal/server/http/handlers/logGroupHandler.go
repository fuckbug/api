package handlers

import (
	"net/http"
	"strconv"

	"github.com/fuckbug/api/internal/middleware"
	logGroup "github.com/fuckbug/api/internal/modules/logGroup"
	"github.com/fuckbug/api/pkg/httputils"
	"github.com/fuckbug/api/pkg/utils"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type logGroupHandler struct {
	logger   Logger
	validate *v.Validate
	service  logGroup.Service
}

func RegisterLogGroupHandlers(
	r *mux.Router,
	logger Logger,
	service logGroup.Service,
	jwtKey []byte,
) {
	h := &logGroupHandler{
		logger:   logger,
		validate: v.New(),
		service:  service,
	}

	routerV1 := r.PathPrefix("/v1/log-groups").Subrouter()
	routerV1.Use(middleware.Auth(jwtKey))

	routerV1.HandleFunc("", h.GetAll).Methods(http.MethodGet)
	routerV1.HandleFunc("/{id}", h.GetByID).Methods(http.MethodGet)
}

// GetByID godoc
// @Summary Get a log group by ID
// @Description Get a log group by ID
// @Tags log-groups
// @Accept json
// @Produce json
// @Success 200 {object} loggroup.Entity
// @Param id path string true "Log ID"
// @Security BearerAuth
// @Router /v1/log-groups/{id} [get].
func (h *logGroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
// @Summary Get all log groups
// @Description Retrieves a list of all logs from the system
// @Tags log-groups
// @Accept  json
// @Produce  json
// @Param projectId query string false "Project ID"
// @Param timeFrom query int false "Time logs from"
// @Param timeTo query int false "Time logs to"
// @Param level query string false "Filter by log level" Enums(DEBUG, INFO, WARN, ERROR)
// @Param search query string false "Search in message field"
// @Param sort query string false "Sort order (asc or desc)" default(desc) Enums(asc, desc)
// @Param limit query int false "Items per page" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} loggroup.EntityList "Successfully retrieved list of logs"
// @Security BearerAuth
// @Router /v1/log-groups [get].
func (h *logGroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	params := logGroup.GetAllParams{
		FilterParams: logGroup.FilterParams{
			ProjectID: projectID,
			TimeFrom:  timeFrom,
			TimeTo:    timeTo,
			Level:     level,
			Search:    search,
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
