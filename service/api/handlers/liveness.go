package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (router *routerImpl) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/* Example of liveness check:
	if err := rt.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/
}
