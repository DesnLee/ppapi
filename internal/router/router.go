package router

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ppapi.desnlee.com/internal/controller"
)

func New() *gin.Engine {
	log.Println("开始初始化路由")

	// 写入日志文件
	// f, _ := os.Create("gin.log")
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(filename)))
	f, _ := os.OpenFile(root+"/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.New()
	_ = r.SetTrustedProxies(nil)

	// 全局使用日志中间件
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s [%s] %s｜%d[%s] \"%s\" <%s %s> %s\n",
			param.TimeStamp.Format(time.DateTime),
			param.ClientIP,
			param.Latency,
			param.StatusCode,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	// 全局使用恢复中间件
	r.Use(gin.Recovery())
	// 全局使用跨域中间件
	// r.Use(middleware.Cors())

	// 静态文件服务
	// r.Static("/static", "./static")

	// 初始化 controllerV1 路由组
	// initV1(r)
	v1 := r.Group("/api/v1")
	v1.GET("/ping", controller.PingHandler)
	v1.POST("/validation_code", controller.SendValidationCodeHandler)

	// 初始化 swagger 路由组
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("初始化路由成功！")
	return r
}
