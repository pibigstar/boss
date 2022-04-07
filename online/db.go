package online

import (
	"encoding/json"
	"fmt"

	"github.com/pibigstar/boss/constant"

	"github.com/pibigstar/boss/config"
	"github.com/pibigstar/boss/core"
	"github.com/pibigstar/boss/logs"
	"github.com/pibigstar/boss/model"
)

// 获取有效的用户
func listAllUser(userId ...int) ([]model.User, error) {
	sql := fmt.Sprintf("select id,username,cookie,status,scoreConfig from user where status = 1")
	if len(userId) > 0 {
		sql += fmt.Sprintf(" and id = %d", userId[0])
	}
	rows, err := config.GetDB().Query(sql)
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
func listAllUserJob() (map[int]map[string]model.Job, error) {
	rows, err := config.GetDB().Query("select userId,jobId,jobName,intervalTime," +
		"spiderTime,helloNum,queueMaxNum,requestResumeTime from jobs where isDel = 0")
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
func addOrUpdateUser(user model.User) error {
	if user.ScoreConfig == "" {
		bs, _ := json.Marshal(core.DefaultScoreConfig())
		user.ScoreConfig = string(bs)
	}
	if user.Status == 0 {
		user.Status = constant.UserStatusActive
		if user.Cookie == "" {
			user.Status = constant.UserStatusDeactivation
		}
	}
	if user.Id == 0 {
		sql := fmt.Sprintf("insert into user(username, cookie, scoreConfig, status) "+
			"values ('%s', '%s', '%s', %d)", user.UserName, user.Cookie, user.ScoreConfig, user.Status)
		_, err := config.GetDB().Exec(sql)
		if err != nil {
			logs.Println("addOrUpdateUser", user, err.Error())
			return err
		}
	} else {
		sql := fmt.Sprintf("update user set username='%s', cookie='%s',scoreConfig='%s', status=%d where id = %d",
			user.UserName, user.Cookie, user.ScoreConfig, user.Status, user.Id)
		_, err := config.GetDB().Exec(sql)
		if err != nil {
			logs.Println("addOrUpdateUser", user, err.Error())
			return err
		}
	}
	return nil
}

// 新增用户Job
func addOrUpdateUserJob(job model.Job) error {
	d := core.DefaultJob(job.JobId, job.JobName)
	if job.HelloNum <= 0 {
		job.HelloNum = d.HelloNum
	}
	if job.SpiderTime <= 0 {
		job.SpiderTime = d.SpiderTime
	}
	if job.IntervalTime <= 0 {
		job.IntervalTime = d.IntervalTime
	}
	if job.QueueMaxNum <= 0 {
		job.QueueMaxNum = d.QueueMaxNum
	}
	if job.RequestResumeTime <= 0 {
		job.RequestResumeTime = d.RequestResumeTime
	}
	if job.Id == 0 {
		sql := fmt.Sprintf("insert into jobs(userId, jobId, jobName, intervalTime, spiderTime, helloNum, queueMaxNum, requestResumeTime)"+
			"values (%d, '%s', '%s', %d, %d, %d, %d, %d)", job.UserId, job.JobId, job.JobName,
			job.IntervalTime, job.SpiderTime, job.HelloNum, job.QueueMaxNum, job.RequestResumeTime)
		_, err := config.GetDB().Exec(sql)
		if err != nil {
			logs.Println("addOrUpdateUserJob", job, err.Error())
			return err
		}
	} else {
		sql := fmt.Sprintf("update jobs set jobName='%s', intervalTime=%d,spiderTime=%d, helloNum=%d, requestResumeTime=%d,isDel=%d where id = %d",
			job.JobName, job.IntervalTime, job.SpiderTime, job.HelloNum, job.RequestResumeTime, job.IsDel, job.Id)
		_, err := config.GetDB().Exec(sql)
		if err != nil {
			logs.Println("addOrUpdateUserJob", job, err.Error())
			return err
		}
	}
	return nil
}

// 显示用户当前Job与可配置的Job

// 获取用户已配置的Job
func listUserJobs(userId int) ([]model.Job, error) {
	sql := fmt.Sprintf("select userId,jobId,jobName,intervalTime,"+
		"spiderTime,helloNum,queueMaxNum,requestResumeTime from jobs where userId = %d and isDel = 0", userId)
	rows, err := config.GetDB().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []model.Job
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
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// 获取用户Boss中设置的Job
func getUser(userId int) (*model.User, error) {
	u := &model.User{}
	sql := fmt.Sprintf("select id,username,cookie,status,scoreConfig from user where id = %d", userId)
	err := config.GetDB().QueryRow(sql).Scan(&u.Id, &u.UserName, &u.Cookie, &u.Status, &u.ScoreConfig)
	if err != nil {
		logs.Println("getUser", err.Error())
		return nil, err
	}
	return u, err
}

// 获取用户Boss中设置的Job
func listUserJobsFromBoss(userId int) ([]*model.Job, error) {
	user, _ := getUser(userId)
	if user == nil {
		return nil, fmt.Errorf("user not exist, id: %d", userId)
	}
	jobs := core.NewBoss(*user).ListJobs()
	return jobs, nil
}

// 重新进行一次招聘
func restart(userId int) {
	if userId == 0 {
		RunOnline()
	} else {
		RunOnline(userId)
	}
}
