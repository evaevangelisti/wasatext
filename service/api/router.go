package api

import (
	"errors"
	"net/http"

	"github.com/evaevangelisti/wasatext/service/api/handlers"
	"github.com/evaevangelisti/wasatext/service/api/middlewares"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/api/services"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger   logrus.FieldLogger
	Database database.Database
}

type Router interface {
	Handler() http.Handler
	Close() error
}

type routerImpl struct {
	httpRouter *httprouter.Router
	logger     logrus.FieldLogger
	database   database.Database
}

func New(config Config) (Router, error) {
	if config.Logger == nil {
		return nil, errors.New("logger is required")
	}

	if config.Database == nil {
		return nil, errors.New("database is required")
	}

	httpRouter := httprouter.New()

	httpRouter.RedirectTrailingSlash = false
	httpRouter.RedirectFixedPath = false

	return &routerImpl{
		httpRouter: httpRouter,
		logger:     config.Logger,
		database:   config.Database,
	}, nil
}

func (router *routerImpl) Handler() http.Handler {
	httpRouter := router.httpRouter

	httpRouter.GET("/liveness", handlers.Liveness(router.database))

	userRepository := &repositories.UserRepository{Database: router.database}
	userService := &services.UserService{Repository: userRepository}
	userHandler := &handlers.UserHandler{Service: userService}

	withAuth := func(handler httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			middlewareHandler := middlewares.AuthMiddleware(userRepository, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r, ps)
			}))

			middlewareHandler.ServeHTTP(w, r)
		}
	}

	httpRouter.GET("/users", withAuth(userHandler.GetUsers))
	httpRouter.GET("/users/:userId", withAuth(userHandler.GetUser))
	httpRouter.POST("/users", userHandler.DoLogin)
	httpRouter.PUT("/me/username", withAuth(userHandler.SetMyUserName))
	httpRouter.PUT("/me/photo", withAuth(userHandler.SetMyPhoto))

	conversationRepository := &repositories.ConversationRepository{Database: router.database}
	conversationService := &services.ConversationService{Repository: conversationRepository}
	conversationHandler := &handlers.ConversationHandler{Service: conversationService}

	httpRouter.GET("/conversations", withAuth(conversationHandler.GetMyConversations))
	httpRouter.GET("/conversations/:conversationId", withAuth(conversationHandler.GetConversation))
	httpRouter.POST("/conversations", withAuth(conversationHandler.CreateConversation))
	httpRouter.POST("/groups/:conversationId/members", withAuth(conversationHandler.AddToGroup))
	httpRouter.PUT("/groups/:conversationId/name", withAuth(conversationHandler.SetGroupName))
	httpRouter.PUT("/groups/:conversationId/photo", withAuth(conversationHandler.SetGroupPhoto))
	httpRouter.DELETE("/groups/:conversationId/members/me", withAuth(conversationHandler.LeaveGroup))

	messageRepository := &repositories.MessageRepository{Database: router.database}
	messageService := &services.MessageService{Repository: messageRepository}
	messageHandler := &handlers.MessageHandler{Service: messageService}

	httpRouter.POST("/conversations/:conversationId/messages", withAuth(messageHandler.SendMessage))
	httpRouter.POST("/conversations/:conversationId/forwards", withAuth(messageHandler.ForwardMessage))
	httpRouter.PUT("/messages/:messageId", withAuth(messageHandler.EditMessage))
	httpRouter.DELETE("/messages/:messageId", withAuth(messageHandler.DeleteMessage))

	commentRepository := &repositories.CommentRepository{Database: router.database}
	commentService := &services.CommentService{Repository: commentRepository}
	commentHandler := &handlers.CommentHandler{Service: commentService}

	httpRouter.POST("/messages/:messageId/comments", withAuth(commentHandler.CommentMessage))
	httpRouter.DELETE("/comments/:commentId", withAuth(commentHandler.UncommentMessage))

	return httpRouter
}

func (router *routerImpl) Close() error {
	return nil
}
