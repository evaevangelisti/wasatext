package handlers

import (
	"net/http"

	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/julienschmidt/httprouter"
)

func Liveness(db database.Database) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
