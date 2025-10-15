package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	SetMessageRoutes(router, db)

	return router
}
