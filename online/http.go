package online

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pibigstar/boss/logs"
)

func RunHttp(port int) {
	r := gin.Default()

	// 修改默认路由打印格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logs.Printf("【%v】====> %v\n", httpMethod, absolutePath)
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
	// 用户列表
	r.GET("/users", func(ctx *gin.Context) {
		users, err := listUser()
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, users)
	})

	// 用户的职位列表
	r.GET("/jobs", func(ctx *gin.Context) {
		userJobs, err := listUserJob()
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, userJobs)
	})

	// 额外信息
	r.GET("/extraInfo", func(ctx *gin.Context) {
		extraInfo := listSchoolAndCompany()
		ctx.JSON(http.StatusOK, extraInfo)
	})
}

func ErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(200, map[string]interface{}{
		"err": err.Error(),
	})
}
