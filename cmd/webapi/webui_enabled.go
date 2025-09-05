//go:build webui

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/evaevangelisti/wasatext/webui"
)

func setupWebUI(handler http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("embedding WebUI dist/ directory: %w", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dashboard" || strings.HasPrefix(r.URL.Path, "/dashboard/") {
			relativePath := strings.TrimPrefix(r.URL.Path, "/dashboard/")
			if relativePath != "" {
				if _, err := fs.Stat(distDirectory, relativePath); err == nil {
					http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
					return
				}
			}

			data, err := fs.ReadFile(distDirectory, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}

			http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(data))
			return
		}

		handler.ServeHTTP(w, r)
	}), nil
}
