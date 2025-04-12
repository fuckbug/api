package server

import (
	"net/http"

	"github.com/fuckbug/api/internal/modules/app"
	"github.com/fuckbug/api/internal/modules/log"
	"github.com/fuckbug/api/internal/server/http/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(
	logger handlers.Logger,
	appService app.Service,
	logService log.Service,
) http.Handler {
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	r.NotFoundHandler = http.HandlerFunc(methodNotFoundHandler)

	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	handlers.RegisterAppHandlers(r, logger, appService)
	handlers.RegisterLogHandlers(r, logger, logService)

	return r
}

func methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func methodNotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
