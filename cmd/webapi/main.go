//go:build webui

/*
This package sets up a web server that serves APIs defined in the `service/api` package.
It connects to necessary external resources, such as a database, and starts two web servers: one for the API and another for debugging purposes.
All API requests are handled by the API web server, while debug variables (/debug/vars) and profiling information (pprof) are served by the debug server.

Usage:

	webapi [flags]

Flags and configurations are automatically managed by the code in `service/api/config/config.go`.

Return values (exit codes):

	0
		The program terminated successfully without errors (e.g., stopped by a signal).

	> 0
		The program terminated due to an error.

Note that this program will automatically update the database schema to the latest version available.
The schema is embedded in the executable during the build process.
*/

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ardanlabs/conf"
	"github.com/evaevangelisti/wasatext/service/api"
	"github.com/evaevangelisti/wasatext/service/config"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/evaevangelisti/wasatext/webui"
	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func run() error {
	rand.Seed(globaltime.Now().UnixNano())

	config, err := config.LoadConfig()

	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}

		return err
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if config.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Infof("initializing application")
	logger.Println("initializing database support")

	db, err := sql.Open("sqlite3", config.Database.FilePath)

	if err != nil {
		logger.WithError(err).Error("failed to open database")
		return fmt.Errorf("opening database: %w", err)
	}

	defer func() {
		logger.Debug("closing database connection")
		db.Close()
	}()

	appDatabase, err := database.New(db)

	if err != nil {
		logger.WithError(err).Error("failed to create AppDatabase instance")
		return fmt.Errorf("creating AppDatabase instance: %w", err)
	}

	logger.Info("initializing API server")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	router, err := api.New(api.Config{
		Logger:   logger,
		Database: appDatabase,
	})

	if err != nil {
		logger.WithError(err).Error("failed to create the API server instance")
		return fmt.Errorf("creating the API server instance: %w", err)
	}

	handler := router.Handler()

	distDirectory, err := fs.Sub(webui.Dist, "dist")

	if err != nil {
		logger.WithError(err).Error("failed to embed WebUI dist/ directory")
		return fmt.Errorf("embedding WebUI dist/ directory: %w", err)
	}

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})

	handler = handlers.CORS(
		handlers.AllowedHeaders([]string{
			"x-example-header",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.MaxAge(1),
	)(handler)

	server := http.Server{
		Addr:              config.Web.APIHost,
		Handler:           handler,
		ReadTimeout:       config.Web.ReadTimeout,
		ReadHeaderTimeout: config.Web.ReadTimeout,
		WriteTimeout:      config.Web.WriteTimeout,
	}

	go func() {
		logger.Infof("API server is listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
		logger.Infof("API server has stopped")
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server encountered an error: %w", err)

	case shutdownSignal := <-shutdown:
		logger.Infof("received signal %v, initiating shutdown", shutdownSignal)

		err := router.Close()

		if err != nil {
			logger.WithError(err).Warning("error during graceful shutdown of API router")
		}

		shutdownContext, cancel := context.WithTimeout(context.Background(), config.Web.ShutdownTimeout)
		defer cancel()

		err = server.Shutdown(shutdownContext)

		if err != nil {
			logger.WithError(err).Warning("error during graceful shutdown of HTTP server")
			err = server.Close()
		}

		switch {
		case shutdownSignal == syscall.SIGSTOP:
			return errors.New("shutdown due to integrity issue")
		case err != nil:
			return fmt.Errorf("could not gracefully stop server: %w", err)
		}
	}

	return nil
}
