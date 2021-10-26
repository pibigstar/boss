package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunHttp(port int) {
	r := gin.Default()

	// 修改默认路由打印格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("【%v】====> %v\n", httpMethod, absolutePath)
	}

	// 使用中间件，阻止panic
	r.Use(gin.Recovery())

	// 注册路由
	InitRouter(r.RouterGroup)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}

func InitRouter(r gin.RouterGroup) {
	r.GET("/users", func(ctx *gin.Context) {
		users, err := listUserFromDB()
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(http.StatusOK, users)
	})
}
