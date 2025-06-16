//go:build !webui

package main

import (
	"net/http"
)

func setupWebUI(handler http.Handler) (http.Handler, error) {
	return handler, nil
}
