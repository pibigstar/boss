package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func sendFeiShu(msg string) {
	uri := "https://open.feishu.cn/open-apis/bot/v2/hook/9b46f934-2e77-499e-81e8-af02b4b27cde"
	text := fmt.Sprintf(`
	{
		"msg_type": "text",
		"content": {
			"text": "%s"
		}
	}`, msg)

	_, err := http.Post(uri, "application/json", strings.NewReader(text))
	if err != nil {
		log.Println("发送飞书提醒失败, msg:", msg)
	}
}


func sendEmail() {
	var (
		username = "741047261@qq.com"
		password = ""
		host     = "smtp.qq.com"
		addr     = "smtp.qq.com:25"
	)
	auth := smtp.PlainAuth("", username, password, host)

	user := "741047261@qq.com"
	to := []string{"741047261@qq.com"}
	msg := []byte(`From: 741047261@qq.com
To: 741047261@qq.com
Subject: Boss登录状态失效

boss登录状态已失效，请及时更改
`)
	err := smtp.SendMail(addr, auth, user, to, msg)
	if err != nil {
		log.Println("发送邮件提醒失败:", err.Error())
	}
}

func isContains(arrs []string, arr string) bool {
	for _, s := range arrs {
		if strings.EqualFold(s, arr) {
			return true
		}
		if strings.Contains(arr, s) {
			return true
		}
		if strings.Contains(s, arr) {
			return true
		}
	}
	return false
}

// 岗位是否匹配
func matchJob(jobName string, geek *Geek) bool {
	jobName = strings.ToLower(jobName)
	expectPositionName := strings.ToLower(geek.GeekCard.ExpectPositionName)
	// 期望职位是否匹配（todo: 这个不准）
	if strings.Contains(jobName, expectPositionName) || strings.Contains(expectPositionName, jobName) {
		return true
	}
	// 个人描述里面是否有该岗位
	if strings.Contains(geek.GeekCard.GeekDesc.Content, jobName) {
		return true
	}
	// 历史工作里面是否有该岗位
	for _, j := range geek.GeekCard.GeekWorks {
		j.PositionName = strings.ToLower(j.PositionName)
		if strings.Contains(j.PositionName, jobName) || strings.Contains(jobName, j.PositionName) {
			return true
		}
	}
	return false
}

func setFilePath() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(currentFile)

	cookieFile = filepath.Join(basePath, cookieFile)
	school985File = filepath.Join(basePath, school985File)
	school211File = filepath.Join(basePath, school211File)
	jobsFile = filepath.Join(basePath, jobsFile)
	companyFile = filepath.Join(basePath, companyFile)
	logFile, _ = os.OpenFile(filepath.Join(basePath, bossLog), os.O_RDWR|os.O_CREATE, 0664)
}

func readCookie() {
	bs, _ := ioutil.ReadFile(cookieFile)
	cookie = string(bs)
}

func readSchool() {
	bs, _ := ioutil.ReadFile(school985File)
	br := bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		school985 = append(school985, string(a))
	}

	bs, _ = ioutil.ReadFile(school211File)
	br = bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		school211 = append(school211, string(a))
	}
}

func readCompany() {
	bs, _ := ioutil.ReadFile(companyFile)
	br := bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		goodCompany = append(goodCompany, string(a))
	}
}

func readJobs() {
	bs, err := ioutil.ReadFile(jobsFile)
	if err != nil {
		return
	}
	br := bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		s := string(a)
		var (
			jobId   string
			jobName string
		)
		if ss := strings.Split(s, "//"); len(s) > 1 {
			jobId = strings.TrimSpace(ss[0])
			jobName = strings.TrimSpace(ss[1])
		}
		if jobId != "" {
			jobIds[jobId] = jobName
		}
	}
}
