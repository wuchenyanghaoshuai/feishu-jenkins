package main

import (
	"fieshu-jenkins/feishu"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	// 计时器
	ticker := time.NewTicker(30 * time.Second)

	// 使用goroutine执行函数，确保定时器不会阻塞主线程
	go func() {
		for {
			select {
			case <-ticker.C:
				feishu.GetInstanceCodeList()
			}
		}
	}()

	r := gin.Default()
	// 设置路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 启动HTTP服务器
	fmt.Println("http://127.0.0.1:8080")
	r.Run(":8080")

}
