package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/router"
)

func Run() {
	rootCmd := &cobra.Command{
		Use: "run",
	}
	svrCmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// 运行服务器
			runServer()
		},
	}
	dbCmd := &cobra.Command{
		Use: "db",
	}
	createCmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			database.CreateTables()
		},
	}
	migrateCmd := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			database.Migrate()
		},
	}
	crudCmd := &cobra.Command{
		Use: "crud",
		Run: func(cmd *cobra.Command, args []string) {
			database.Crud()
		},
	}

	rootCmd.AddCommand(svrCmd, dbCmd)
	dbCmd.AddCommand(createCmd, migrateCmd, crudCmd)

	// 连接数据库
	defer database.Close()
	database.Connect()

	// 命令行运行
	err := rootCmd.Execute()
	if err != nil {
		return
	}

}

func runServer() {

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
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
