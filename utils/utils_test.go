package utils

import (
	"github.com/pibigstar/boss"
	"testing"
)

func TestReadSchool(t *testing.T) {
	readSchool()
	for _, s := range main.school985 {
		t.Log(s)
	}

	for _, s := range main.school211 {
		t.Log(s)
	}
}

func TestReadJob(t *testing.T) {
	readJobs()
	for jobId, jobName := range main.jobIds {
		t.Log(jobId, jobName)
	}
}

func TestReadCompany(t *testing.T) {
	readCompany()
	for _, c := range main.goodCompany {
		t.Log(c)
	}
}

func TestSendEmail(t *testing.T) {
	SendEmail()
}

func TestSendFeiShu(t *testing.T) {
	SendFeiShu("Boss当前登录登录状态失效")
}