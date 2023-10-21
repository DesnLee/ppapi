package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"ppapi.desnlee.com/docs"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/email"
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
			port := "9999"
			if len(args) != 0 {
				port = args[0]
			}
			runServer(port)
		},
	}

	emailCmd := &cobra.Command{
		Use: "email",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("请输入收件人邮箱列表")
			}
			email.Send(args)
		},
	}

	coverCmd := &cobra.Command{
		Use: "cover",
		Run: func(cmd *cobra.Command, args []string) {
			if err := exec.Command("MailHog").Start(); err != nil {
				log.Fatalln("MailHog 启动失败：", err)
			}
			if err := exec.Command("go", "test", "-coverprofile=coverage.out", "./...").Run(); err != nil {
				log.Fatalln("测试失败：", err)
			}
			if err := exec.Command("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html").Run(); err != nil {
				log.Fatalln("生成 html 失败：", err)
			}

			var port string
			if len(args) == 0 {
				port = "8888"
			} else {
				port = args[0]
			}

			defer func() {
				os.Remove("coverage.out")
				os.Remove("coverage.html")
			}()
			log.Printf("查看覆盖率报告 http://localhost:%v/coverage.html", port)
			if err := http.ListenAndServe(":"+port, http.FileServer(http.Dir("."))); err != nil {
				log.Fatalln("服务器启动失败：", err)
			}
		},
	}

	dbCmd := &cobra.Command{
		Use: "db",
	}
	newMigrationCmd := &cobra.Command{
		Use: "migrate:new",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("请输入迁移名称")
			}
			database.MigrateNew(args[0])
		},
	}
	migrateUpCmd := &cobra.Command{
		Use: "migrate:up",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateUp()
		},
	}
	migrateDownCmd := &cobra.Command{
		Use: "migrate:down",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("请输入回退步数")
			}
			step, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			database.MigrateDown(step)
		},
	}
	crudCmd := &cobra.Command{
		Use: "crud",
		Run: func(cmd *cobra.Command, args []string) {
			database.Crud()
		},
	}

	rootCmd.AddCommand(svrCmd, emailCmd, coverCmd, dbCmd)
	dbCmd.AddCommand(newMigrationCmd, migrateUpCmd, migrateDownCmd, crudCmd)

	// 连接数据库
	defer database.Close()
	database.Connect()

	// 命令行运行
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

func runServer(port string) {
	// 初始化服务器
	docs.SwaggerInfo.Host = "localhost:" + port

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
