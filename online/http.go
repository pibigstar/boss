package online

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pibigstar/boss/model"

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
		users, err := listAllUser()
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, users)
	})

	// 用户的职位列表
	r.GET("/jobs", func(ctx *gin.Context) {
		userJobs, err := listAllUserJob()
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

	// 新增修改用户
	r.POST("/addOrUpdateUser", func(ctx *gin.Context) {
		var user model.User
		if err := ctx.ShouldBind(&user); err != nil {
			ErrorResponse(ctx, err)
			return
		}
		err := addOrUpdateUser(user)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, "添加成功")
	})

	// 新增修改用户招聘岗位
	r.POST("/addOrUpdateUserJob", func(ctx *gin.Context) {
		var job model.Job
		if err := ctx.ShouldBind(&job); err != nil {
			ErrorResponse(ctx, err)
			return
		}
		err := addOrUpdateUserJob(job)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, "添加成功")
	})

	// 获取用户已配置的职位列表
	r.GET("/listUserJob", func(ctx *gin.Context) {
		idStr := ctx.Query("userId")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		jobs, err := listUserJobs(userId)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, jobs)
	})

	// 从Boss中获取用户配置的职位列表
	r.GET("/listUserJobFromBoss", func(ctx *gin.Context) {
		idStr := ctx.Query("userId")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		jobs, err := listUserJobsFromBoss(userId)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, jobs)
	})

	// 重新进行一次招聘
	r.GET("/restart", func(ctx *gin.Context) {
		idStr := ctx.Query("userId")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorResponse(ctx, err)
			return
		}
		restart(userId)
		ctx.JSON(http.StatusOK, "重启成功")
	})
}

func ErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(200, map[string]interface{}{
		"err": err.Error(),
	})
}
