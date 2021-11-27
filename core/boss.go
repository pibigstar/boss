package core

import "sync"

type Bosser interface {
	Hiring() error
}

type Boss struct {
	talked sync.Map
	Cookie string
	Jobs   map[string]Job
}

type Job struct {
	JobId             string `json:"jobId"`
	JobName           string `json:"jobName"`
	IntervalTime      int    `json:"intervalTime"`      // 间隔时间，多久抓一次简历列表，单位s默认10s
	SpiderTime        int    `json:"spiderTime"`        // 抓取时长，单位 s，默认 180s
	HelloNum          int    `json:"helloNum"`          // 对几个人打招呼
	QueueMaxNum       int    `json:"queueMaxNum"`       // 最多可以有多少个候选人
	RequestResumeTime int    `json:"requestResumeTime"` // 请求简历时长,单位 hour, 默认 1h
}

func DefaultJob(jobId, jobName string) Job {
	j := Job{
		IntervalTime:      10,
		SpiderTime:        180,
		HelloNum:          3,
		QueueMaxNum:       10,
		RequestResumeTime: 1,
	}
	return j
}
