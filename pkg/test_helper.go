package pkg

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "ppapi.desnlee.com/config"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/jwt_helper"
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

func TestCreateUserAndJWT() (sqlcExec.User, string) {
	// 创建用户
	u, err := database.Q.CreateUser(database.DBCtx, "test@qq.com")
	if err != nil {
		log.Fatal("Create User Error: ", err)
	}

	// 生成 JWT
	jwtStr, err := jwt_helper.GenerateJWT(u.ID)
	if err != nil {
		log.Fatal("Generate JWT Error: ", err)
	}

	return u, jwtStr
}
