package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/evaevangelisti/wasatext/service/api/middlewares"
	"github.com/evaevangelisti/wasatext/service/api/services"
	"github.com/evaevangelisti/wasatext/service/utils"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	Service *services.UserService
}

type GetUsersQuery struct {
	Q string `validate:"omitempty,min=1,max=16"`
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

	query := GetUsersQuery{Q: r.URL.Query().Get("q")}

	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	users, err := handler.Service.GetUsers(query.Q, auid)
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

	user, created, err := handler.Service.DoLogin(request.Username)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

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

	var dstFilename, dstPath string

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != utils.ExtJPG && ext != utils.ExtJPEG && ext != utils.ExtPNG && ext != utils.ExtWEBP {
			errors.WriteHTTPError(w, errors.ErrBadRequest)
			return
		}

		dstFilename = uuid.New().String() + ext
		dstPath = "./tmp/uploads/profile-pictures/" + dstFilename
	}

	profilePicture := ""
	if dstFilename != "" {
		profilePicture = "/uploads/profile-pictures/" + dstFilename
	}

	user, err := handler.Service.UpdateProfilePicture(auid, profilePicture)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	if file != nil && dstPath != "" {
		dstDir := "./tmp/uploads/profile-pictures"
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
			return
		}

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
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}
