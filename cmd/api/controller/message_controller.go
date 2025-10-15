package controller

import (
	"chat_service/internal/domain/message"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController interface {
	GetPaginated(ctx *gin.Context)
}

type messageController struct {
	repository message.MessageRepository
}

func NewMessageController(repository message.MessageRepository) MessageController {
	return &messageController{repository}
}

func (c *messageController) GetPaginated(ctx *gin.Context) {
	idChat := ctx.Param("idChat")
	offset := ctx.Param("offset")

	parsedOffset, _ := strconv.ParseInt(offset, 10, 64)

	messages, err := c.repository.GetPaginated(idChat, int(parsedOffset))
	if err != nil {
		ctx.JSON(500, gin.H{"message": err})
		return
	}

	ctx.JSON(200, messages)
}
