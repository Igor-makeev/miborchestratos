package handler

import (
	"miborchestrator/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
	Router  *gin.Engine
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		Service: service,
	}
	auth := handler.Router.Group("/auth")
	{
		auth.POST("/register", handler.register)
		auth.POST("/login", handler.login)
	}
	api := handler.Router.Group("/api", handler.useridentity)
	{
		user := api.Group("/user")
		{
			user.POST("/create_wallet", handler.sendToCreateQueue)
			user.POST("/do_transaction", handler.initTransfer)

		}

		api.POST("/send_to_txManager", handler.sendToTxManager)
		//Todo: добавить в ридми описание хендлера
	}

	return handler
}
