package route

import (
	"chat_service/cmd/api/controller"
	"chat_service/cmd/api/middleware"
	"chat_service/internal/domain/message"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetMessageRoutes(router *gin.Engine, db *gorm.DB) {
	messageRepository := message.NewMessageRepository(db)
	messageController := controller.NewMessageController(messageRepository)

	router.GET("/message/:idChat/:offset", middleware.ProtectRoute(), messageController.GetPaginated)
}
