package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupStaticRoutes(router *gin.Engine) {
	// 设置静态文件夹
	router.Static("/assets", "./dist/assets") // 这里的 "./static" 是 Vue 打包后的静态文件所在的目录

	// 处理 Vue 路由入口 HTML 文件
	router.GET("/", func(c *gin.Context) {
		c.File("./dist/index.html") // 这里的 "./static/index.html" 是 Vue 应用程序的入口 HTML 文件路径
	})

	// 处理其他所有未被匹配的 GET 请求
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.File("./dist/index.html")
		}
	})
}
