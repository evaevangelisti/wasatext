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

type MessageHandler struct {
	Service *services.MessageService
}

type SendMessageRequest struct {
	Content string `validate:"omitempty,min=1,max=1000"`
}

func (handler *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	content := r.FormValue("content")

	request := SendMessageRequest{content}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var attachment string

	file, header, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext == "" {
			errors.WriteHTTPError(w, errors.ErrBadRequest)
			return
		}

		dstDir := "./uploads/attachments"
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

		attachment = "/uploads/attachments/" + dstFilename
	}

	message, err := handler.Service.CreateMessage(cid, auid, content, attachment)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(message); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

type ForwardMessageRequest struct {
	MessageID uuid.UUID `json:"messageId" validate:"required,uuid4"`
}

func (handler *MessageHandler) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var request ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	forwardedMessage, err := handler.Service.CreateForwardedMessage(cid, auid, request.MessageID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(forwardedMessage); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

type EditMessageRequest struct {
	Content string `json:"content,omitempty" validate:"omitempty,min=1,max=1000"`
}

func (handler *MessageHandler) EditMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	messageID := ps.ByName("messageId")

	mid, err := uuid.Parse(messageID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var request EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	updatedMessage, err := handler.Service.UpdateMessage(mid, auid, request.Content)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(updatedMessage); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	messageID := ps.ByName("messageId")

	mid, err := uuid.Parse(messageID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	err = handler.Service.DeleteMessage(mid, auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
