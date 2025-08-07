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

type ConversationHandler struct {
	Service *services.ConversationService
}

func (handler *ConversationHandler) GetMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	conversations, err := handler.Service.GetConversationsByUserID(auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(conversations); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	conversationID := ps.ByName("conversationId")

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	conversation, err := handler.Service.GetConversationByID(cid, auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	if conversation == nil {
		errors.WriteHTTPError(w, errors.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(conversation); err != nil {
		return
	}
}

type CreateConversationRequest struct {
	Type    string      `json:"type" validate:"required,oneof=private group"`
	UserID  uuid.UUID   `json:"userId,omitempty" validate:"omitempty"`
	Name    string      `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Members []uuid.UUID `json:"members,omitempty" validate:"omitempty,max=100"`
}

func (handler *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	var request CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	if request.UserID != uuid.Nil && len(request.Members) > 0 {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	switch request.Type {
	case "private":
		if request.UserID == uuid.Nil || request.UserID == auid {
			errors.WriteHTTPError(w, errors.ErrBadRequest)
			return
		}

		participants := []uuid.UUID{auid, request.UserID}

		privateConversation, err := handler.Service.CreatePrivateConversation(participants)
		if err != nil {
			if err == errors.ErrConflict {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(errors.ErrConflict.StatusCode)

				if err = json.NewEncoder(w).Encode(privateConversation); err != nil {
					errors.WriteHTTPError(w, errors.ErrInternal)
				}

				return
			}

			errors.WriteHTTPError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(privateConversation); err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
		}
	case "group":
		members := append([]uuid.UUID{auid}, request.Members...)

		groupConversation, err := handler.Service.CreateGroupConversation(request.Name, members)
		if err != nil {
			errors.WriteHTTPError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(groupConversation); err != nil {
			errors.WriteHTTPError(w, errors.ErrInternal)
		}
	default:
		errors.WriteHTTPError(w, errors.ErrBadRequest)
	}
}

type AddToGroupRequest struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
}

func (handler *ConversationHandler) AddToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	conversationID := ps.ByName("conversationId")

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var request AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	groupConversation, err := handler.Service.AddMember(cid, auid, request.UserID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(groupConversation); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

type SetGroupNameRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

func (handler *ConversationHandler) SetGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	conversationID := ps.ByName("conversationId")

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var request SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	groupConversation, err := handler.Service.UpdateGroupName(cid, auid, request.Name)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(groupConversation); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *ConversationHandler) SetGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	conversationID := ps.ByName("conversationId")

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
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
		dstPath = "./tmp/uploads/group-photos/" + dstFilename
	}

	photo := ""
	if dstFilename != "" {
		photo = "/uploads/group-photos/" + dstFilename
	}

	groupConversation, err := handler.Service.UpdateGroupPhoto(cid, auid, photo)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	if file != nil && dstPath != "" {
		dstDir := "./tmp/uploads/group-photos/"
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

	if err = json.NewEncoder(w).Encode(groupConversation); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *ConversationHandler) LeaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	conversationID := ps.ByName("conversationId")

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	err = handler.Service.RemoveMember(cid, auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
