package online

import "github.com/pibigstar/boss/core"

// 在线招聘主函数

func RunOnline() {

	// 1. 读取可用用户
	users, err := listUser()
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		return
	}

	// 2. 读取211学校、985学校、大厂信息
	sc := listSchoolAndCompany()

	// 4. 读取用户配置的需要招聘的job

	userJobs, err := listUserJob()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		if jobs, ok := userJobs[user.Id]; ok {
			boss := core.NewBoss(user, jobs, sc)
			go func() {
				defer func() {
					recover()
				}()
				boss.Hiring()
			}()
		}
	}

	// 每个job抓取几分钟，对几个人进行打招呼
	// 最大的候选人队列值

	// 5. for 对 job 进行招聘
}
