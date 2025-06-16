//go:build webui

package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/evaevangelisti/wasatext/webui"
)

func setupWebUI(handler http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")

	if err != nil {
		return nil, fmt.Errorf("embedding WebUI dist/ directory: %w", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	}), nil
}
