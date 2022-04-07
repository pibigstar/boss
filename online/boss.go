package online

import (
	"encoding/json"
	"log"

	"github.com/pibigstar/boss/core"
	"github.com/robfig/cron"
)

// 在线招聘主函数

func RunOnline(userId ...int) {
	// 1. 读取可用用户
	users, err := listAllUser(userId...)
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		return
	}

	// 2. 读取211学校、985学校、大厂信息
	sc := listSchoolAndCompany()

	// 4. 读取用户配置的需要招聘的job
	userJobs, err := listAllUserJob()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		if jobs, ok := userJobs[user.Id]; ok {
			boss := core.NewBoss(user, func(boss *core.Boss) {
				boss.Jobs = jobs
				boss.ExtraInfo = sc
				if user.ScoreConfig != "" {
					err := json.Unmarshal([]byte(user.ScoreConfig), &boss.ScoreConfig)
					if err != nil {
						log.Println("unmarshal score config", user.ScoreConfig)
					}
				}
			})
			go func() {
				defer func() {
					recover()
				}()
				// 5. 开始进行招聘
				boss.Hiring()
			}()
		}
	}
}

func CronRun() {
	c := cron.New()
	// 早上九点半开始执行
	err := c.AddFunc("0 30 9 * * *", func() {
		RunOnline()
	})
	if err != nil {
		panic(err)
	}
}
