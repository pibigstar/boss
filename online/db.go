package online

import (
	"github.com/pibigstar/boss/config"
	"github.com/pibigstar/boss/logs"
)

// 获取有效的用户
func listUser() ([]*User, error) {
	rows, err := config.GetDB().Query("select id,username,cookie,status from user where status = 1")
	if err != nil {
		return nil, err
	}
	var users []*User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.UserName, &u.Cookie, &u.Status)
		if err != nil {
			logs.Println("list user rows scan", err.Error())
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func listUserJob() (map[int][]*Job, error) {
	rows, err := config.GetDB().Query("select userId,jobId,jobName,intervalTime," +
		"spiderTime,helloNum,queueMaxNum,requestResumeTime from jobs where status = 1")
	if err != nil {
		return nil, err
	}
	userJob := make(map[int][]*Job)
	for rows.Next() {
		var j Job
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
		userJob[j.UserId] = append(userJob[j.UserId], &j)
	}
	return userJob, nil
}

func GetUserConfig() {

}
