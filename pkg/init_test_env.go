package pkg

import (
	"github.com/gin-gonic/gin"
	_ "ppapi.desnlee.com/config"
	"ppapi.desnlee.com/internal/database"
)

func InitTestEnv() (*gin.Engine, func()) {
	database.Connect()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	return r, func() {
		database.Close()
	}
}
