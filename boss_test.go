package main

import (
	"context"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

// 测试招人
func TestHiring(t *testing.T) {
	for jobId, jobName := range jobIds {
		Hiring(jobId, jobName)
	}
}

func TestListRecommend(t *testing.T) {
	for jobId := range jobIds {
		geeks, err := listRecommend(jobId)
		if err != nil {
			t.Error(err)
		}
		for _, geek := range geeks {
			t.Log(geek.GeekCard.GeekName)
		}
	}
}

func TestSetHelloMsg(t *testing.T) {
	setHelloMsg()
}

func TestGetQR(t *testing.T) {
	getQRId(context.Background())
}

func TestListJob(t *testing.T) {
	jobs := listJobs()
	for _, j := range jobs {
		t.Log(j.JobName, j.JobId)
	}
}

func TestInputJobs(t *testing.T) {
	inputJobs()
}
