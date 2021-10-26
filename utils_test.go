package main

import "testing"

func TestReadSchool(t *testing.T) {
	readSchool()
	for _, s := range school985 {
		t.Log(s)
	}

	for _, s := range school211 {
		t.Log(s)
	}
}

func TestReadJob(t *testing.T) {
	readJobs()
	for jobId, jobName := range jobIds {
		t.Log(jobId, jobName)
	}
}

func TestReadCompany(t *testing.T) {
	readCompany()
	for _, c := range goodCompany {
		t.Log(c)
	}
}

func TestSendEmail(t *testing.T) {
	sendEmail()
}

func TestSendFeiShu(t *testing.T) {
	sendFeiShu("Boss当前登录登录状态失效")
}