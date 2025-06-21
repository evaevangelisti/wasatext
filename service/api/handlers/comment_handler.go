package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/evaevangelisti/wasatext/service/api/middlewares"
	"github.com/evaevangelisti/wasatext/service/api/services"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type CommentHandler struct {
	Service *services.CommentService
}

type CommentMessageRequest struct {
	Emoji string `json:"emoji" validate:"required"`
}

func (handler *CommentHandler) CommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var request CommentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	re := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)
	if !re.MatchString(request.Emoji) {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	comment, err := handler.Service.CreateComment(mid, auid, request.Emoji)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(comment); err != nil {
		errors.WriteHTTPError(w, errors.ErrInternal)
	}
}

func (handler *CommentHandler) UncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	commentID := ps.ByName("commentId")

	cid, err := uuid.Parse(commentID)
	if err != nil {
		errors.WriteHTTPError(w, errors.ErrBadRequest)
		return
	}

	err = handler.Service.DeleteComment(cid, auid)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
