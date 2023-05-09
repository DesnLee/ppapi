package main

import (
	"log"

	"ppapi.desnlee.com/cmd"
)

func main() {
	defer log.Println("服务已关闭")

	// 初始化服务器
	cmd.RunServer()
}
