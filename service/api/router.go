package api

import (
	"errors"
	"net/http"

	"github.com/evaevangelisti/wasaphoto/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger   logrus.FieldLogger
	Database database.Database
}

type Router interface {
	Handler() http.Handler
	Close() error
}

type routerImpl struct {
	router   *httprouter.Router
	logger   logrus.FieldLogger
	database database.Database
}

func New(config Config) (Router, error) {
	if config.Logger == nil {
		return nil, errors.New("logger is required")
	}

	if config.Database == nil {
		return nil, errors.New("database is required")
	}

	router := httprouter.New()

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &routerImpl{
		router:   router,
		logger:   config.Logger,
		database: config.Database,
	}, nil
}

func (router *routerImpl) Handler() http.Handler {
	return router.router
}

func (router *routerImpl) Close() error {
	return nil
}
