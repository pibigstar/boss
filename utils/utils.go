package utils

import (
	"fmt"
	"github.com/pibigstar/boss/logs"
	"github.com/pibigstar/boss/model"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
)

func SendFeiShu(msg string) {
	logs.Println(msg)
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
		logs.Println("发送飞书提醒失败, msg:", msg)
	}
}

func SendEmail() {
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
		logs.Println("发送邮件提醒失败:", err.Error())
	}
}

// 岗位是否匹配
func MatchJob(jobName string, geek *model.Geek) bool {
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

func IsContains(arrs []string, arr string) bool {
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

func MapLen(m *sync.Map) int {
	var i int
	m.Range(func(key, value interface{}) bool {
		i++
		return true
	})
	return i
}
