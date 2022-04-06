package online

import (
	"github.com/pibigstar/boss/config"
	"github.com/pibigstar/boss/logs"
	"github.com/pibigstar/boss/model"
)

// 获取有效的用户
func listUser() ([]model.User, error) {
	rows, err := config.GetDB().Query("select id,username,cookie,status,scoreConfig from user where status = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.Id, &u.UserName, &u.Cookie, &u.Status, &u.ScoreConfig)
		if err != nil {
			logs.Println("list user rows scan", err.Error())
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// 获取用户当前的招聘岗位配置
func listUserJob() (map[int]map[string]model.Job, error) {
	rows, err := config.GetDB().Query("select userId,jobId,jobName,intervalTime," +
		"spiderTime,helloNum,queueMaxNum,requestResumeTime from jobs where status = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userJob := make(map[int]map[string]model.Job)
	for rows.Next() {
		var j model.Job
		err = rows.Scan(
			&j.UserId,
			&j.JobId,
			&j.JobName,
			&j.IntervalTime,
			&j.SpiderTime,
			&j.HelloNum,
			&j.QueueMaxNum,
			&j.RequestResumeTime)
		if err != nil {
			logs.Println("list user rows scan", err.Error())
			continue
		}
		if r, ok := userJob[j.UserId]; ok {
			r[j.JobName] = j
			userJob[j.UserId] = r
		} else {
			r = make(map[string]model.Job)
			r[j.JobName] = j
			userJob[j.UserId] = r
		}
	}
	return userJob, nil
}

// 获取学校与公司相关信息
func listSchoolAndCompany() map[string][]string {
	rows, err := config.GetDB().Query("select name,tag from extraInfo")
	if err != nil {
		return nil
	}
	defer rows.Close()
	sc := make(map[string][]string)
	for rows.Next() {
		var e model.ExtraInfo
		err = rows.Scan(&e.Name, &e.Tag)
		if err != nil {
			logs.Println("listSchoolAndCompany", err.Error())
			continue
		}
		sc[e.Tag] = append(sc[e.Tag], e.Name)
	}
	return sc
}

// 新增用户

// 修改用户打分配置

// 修改用户cookie

// 新增用户Job

// 显示用户当前Job与可配置的Job
