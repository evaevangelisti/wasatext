package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/evaevangelisti/wasatext/service/api/middlewares"
	"github.com/evaevangelisti/wasatext/service/api/services"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	Service *services.UserService
}

func (handler *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authenticatedUserID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	auid, err := uuid.Parse(authenticatedUserID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")

	users, err := handler.Service.GetUsers(q, auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(users); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userId")

	uid, err := uuid.Parse(userID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	user, err := handler.Service.GetUserByID(uid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	if user == nil {
		errors.WriteHTTPError(w, errors.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

type DoLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
}

func (handler *UserHandler) DoLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request DoLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	user, err := handler.Service.DoLogin(request.Username)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

type SetMyUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
}

func (handler *UserHandler) SetMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authenticatedUserID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	auid, err := uuid.Parse(authenticatedUserID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	var request SetMyUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	user, err := handler.Service.UpdateUsername(auid, request.Username)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *UserHandler) SetMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authenticatedUserID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	auid, err := uuid.Parse(authenticatedUserID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var profilePicture string

	file, header, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext == "" {
			errors.WriteHTTPError(w, errors.ErrBadRequest)
			return
		}

		dstDir := "./uploads/profile-pictures"
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
			return
		}

		dstFilename := uuid.New().String() + ext
		dstPath := filepath.Join(dstDir, dstFilename)

		dst, err := os.Create(dstPath)
		if err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
			return
		}

		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
			return
		}

		profilePicture = "/uploads/profile-pictures/" + dstFilename
	}

	user, err := handler.Service.UpdateProfilePicture(auid, profilePicture)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}
