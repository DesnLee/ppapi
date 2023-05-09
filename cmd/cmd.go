package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/router"
)

func RunServer() {
	// 连接数据库
	defer database.Close()
	database.Connect()
	database.CreateTables()

	// 初始化服务器
	const port = "9999"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.New(),
	}

	// 启动服务器
	go func() {
		// 服务连接
		log.Printf("服务器已启动于 http://localhost:%v\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败：%s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("开始关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("关闭失败：", err)
	}
}
