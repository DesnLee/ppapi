package pkg

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "ppapi.desnlee.com/config"
	"ppapi.desnlee.com/internal/database"
)

func InitTestEnv() (*gin.Engine, func()) {
	database.Connect()

	if err := database.Q.DeleteAllUser(database.DBCtx); err != nil {
		log.Fatalln("Delete All User Error: ", err)
	}

	if err := database.Q.DeleteAllValidationCode(database.DBCtx); err != nil {
		log.Fatalln("Delete All Validation Code Error: ", err)
	}

	if err := database.Q.DeleteAllTag(database.DBCtx); err != nil {
		log.Fatalln("Delete All Tag Error: ", err)
	}

	if err := database.Q.DeleteAllItem(database.DBCtx); err != nil {
		log.Fatalln("Delete All Item Error: ", err)
	}

	if err := database.Q.DeleteAllItemTagRelation(database.DBCtx); err != nil {
		log.Fatalln("Delete All Item Tag Relation Error: ", err)
	}

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	return r, func() {
		database.Close()
	}
}
