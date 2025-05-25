package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/fuckbug/api/docs" // for swagger

	"github.com/fuckbug/api/internal/modules/app"
	"github.com/fuckbug/api/internal/storage/sql"

	"github.com/fuckbug/api/internal/logger"
	moduleError "github.com/fuckbug/api/internal/modules/errors"
	moduleGroupError "github.com/fuckbug/api/internal/modules/errorsGroup"
	moduleLog "github.com/fuckbug/api/internal/modules/log"
	moduleGroupLog "github.com/fuckbug/api/internal/modules/logGroup"
	moduleProject "github.com/fuckbug/api/internal/modules/project"
	moduleUser "github.com/fuckbug/api/internal/modules/users"
	server "github.com/fuckbug/api/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/fuckbug/config.json", "Path to configuration file")
}

const serverShutdownTimeout = 3 * time.Second

// @title FuckBug API
// @version 1.0.0
// @description This is FuckBug.io API.
// @termsOfService https://fuckbug.io/terms/
// @contact.name API Support
// @contact.email support@fuckbug.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	flag.Parse()

	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	appLogger := logger.New(config.Logger.Level, nil)

	db, err := sqlx.Connect("postgres", config.Postgres.Dsn)
	if err != nil {
		appLogger.Error(fmt.Sprintf("failed to connect to database: %v", err))
		return
	}

	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	if err := sql.RunMigrations(db, appLogger); err != nil {
		appLogger.Error(fmt.Sprintf("failed to run migrations: %v", err))
		return
	}

	jwtKey := []byte("teststringjwt") // todo

	appService := app.New(appLogger)
	userService := moduleUser.NewService(moduleUser.NewRepository(db, appLogger), jwtKey, appLogger)
	logService := moduleLog.NewService(moduleLog.NewRepository(db, appLogger), appLogger)
	logGroupService := moduleGroupLog.NewService(moduleGroupLog.NewRepository(db, appLogger), appLogger)
	errorService := moduleError.NewService(moduleError.NewRepository(db, appLogger), appLogger)
	errorGroupService := moduleGroupError.NewService(moduleGroupError.NewRepository(db, appLogger), appLogger)
	projectService := moduleProject.NewService(moduleProject.NewRepository(db, appLogger), appLogger, config.Domain)

	s := server.New(
		appLogger,
		appService,
		userService,
		logService,
		logGroupService,
		errorService,
		errorGroupService,
		projectService,
		"",
		config.Port,
		jwtKey,
	)

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, serverShutdownTimeout)
		defer shutdownCancel()

		if err := s.Stop(shutdownCtx); err != nil {
			appLogger.Error("failed to stop http server: " + err.Error())
		}
	}()

	appLogger.Info(fmt.Sprintf("Service listening on port: %d", config.Port))

	if err := s.Start(ctx); err != nil {
		appLogger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
