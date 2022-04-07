package core

import (
	"sync"

	"github.com/pibigstar/boss/model"
)

type Boss struct {
	talked      sync.Map
	Cookie      string
	User        model.User           // 当前用户相关信息（仅在线情况下有值）
	Jobs        map[string]model.Job // 招聘岗位相关配置
	ScoreConfig model.ScoreConfig    // 候选人筛选打分配置
	ExtraInfo   map[string][]string  // 985,211,大厂信息
}

type Option func(boss *Boss)

func NewBoss(user model.User, options ...Option) *Boss {
	boss := &Boss{
		Cookie:      user.Cookie,
		User:        user,
		ScoreConfig: DefaultScoreConfig(),
	}
	for _, o := range options {
		o(boss)
	}
	return boss
}

// DefaultJob 基础的Job配置信息
func DefaultJob(jobId, jobName string) model.Job {
	return model.Job{
		JobId:             jobId,
		JobName:           jobName,
		IntervalTime:      10,
		SpiderTime:        180,
		HelloNum:          3,
		QueueMaxNum:       10,
		RequestResumeTime: 1,
	}
}

// DefaultScoreConfig 最基础的分值配置信息
func DefaultScoreConfig() model.ScoreConfig {
	return model.ScoreConfig{
		Score211:     3,
		Score958:     4,
		GoodCompany:  5,
		Undergrad:    3,
		Master:       4,
		WorkTime:     3,
		AgeOver35:    -4,
		OnlineWork:   2,
		OfflineWork:  3,
		ActiveToday:  1,
		ActiveMinute: 2,
	}
}
