package pkg

import (
	"github.com/gin-gonic/gin"
	_ "ppapi.desnlee.com/config"
	"ppapi.desnlee.com/internal/database"
)

func SetupTest() *gin.Engine {
	database.Connect()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	return r
}
