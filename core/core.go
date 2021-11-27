package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pibigstar/boss/config"
	"github.com/pibigstar/boss/constant"
	"github.com/pibigstar/boss/logs"
	"github.com/pibigstar/boss/model"
	"github.com/pibigstar/boss/utils"
)

var (
	maxLimit       = errors.New("今日沟通已达上限")
	ErrTalked      = errors.New("候选人已沟通过")
	notFriend      = errors.New("好友关系校验失败")
	notLogin       = errors.New("当前登录状态已失效")
	accountUnusual = errors.New("账号异常")
)

var (
	client = http.Client{
		Timeout: 5 * time.Second,
	}
)

// 招人
func (boss *Boss) Hiring() {
	var (
		wg sync.WaitGroup
	)

	var hiring = func(job Job) {
		// 10秒一次，防止被反爬
		t := time.NewTicker(time.Duration(job.IntervalTime) * time.Second)
		// 进行3分钟的候选人选择
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.SpiderTime)*time.Second)

		defer func() {
			t.Stop()
			cancel()
			if e := recover(); e != nil {
				logs.Println("hiring recover", e)
			}
		}()

		var geeksQueue []*model.Geek

		for {
			select {
			case <-t.C:
				geeks, err := boss.searchGeekByJobId(job.JobId, job.JobName)
				if err != nil {
					if err == notLogin {
						utils.SendFeiShu("Boss当前登录状态失效")
					}
					// 通知可以去打招呼了
					cancel()
					t.Stop()
				}
				geeksQueue = append(geeksQueue, geeks...)
				// 候选人达到最大值就没必要继续跑了
				if len(geeksQueue) == job.QueueMaxNum {
					cancel()
					t.Stop()
				}

			case <-ctx.Done():
				// 打招呼并请求简历
				boss.helloAndRequestResumes(job.JobId, job.RequestResumeTime, geeksQueue)
				return
			}
		}
	}

	for _, job := range boss.Jobs {
		wg.Add(1)
		go func(job Job) {
			defer wg.Done()
			hiring(job)
		}(job)
	}
	wg.Wait()
}

// 打招呼并轮询请求简历
func (boss *Boss) helloAndRequestResumes(jobId string, requestResumeTime int, geeksQueue []*model.Geek) {
	// 按权重排序
	sort.Sort(model.SortGeek(geeksQueue))
	var wg sync.WaitGroup
	// 进行一小时的请求简历
	rrCtx, cancel := context.WithTimeout(context.Background(), time.Duration(requestResumeTime)*time.Hour)
	defer cancel()

	for _, l := range geeksQueue {
		err := boss.hello(jobId, l.GeekCard.GeekID, l.GeekCard.EncryptGeekID, l.GeekCard.Lid, l.GeekCard.SecurityID, l.GeekCard.ExpectID)
		if err != nil {
			if err == maxLimit {
				logs.Println("已达上限")
				break
			}
			logs.Printf("打招呼失败，候选人: %s, 分值: %d err: %s\n", l.GeekCard.GeekName, l.Weight, err.Error())
			continue
		}
		msg := fmt.Sprintf("正在与: %s 打招呼, 分值: %d", l.GeekCard.GeekName, l.Weight)
		utils.SendFeiShu(msg)

		// 轮询向牛人直接请求简历直到对方回复我们建立好友关系
		wg.Add(1)
		go func(name, securityId string) {
			defer wg.Done()
			t := time.NewTicker(time.Minute * 5)
			for {
				select {
				case <-t.C:
					logs.Printf("正在索求候选人:%s的简历 \n", name)
					if err := boss.requestResumes(name, securityId); err == nil {
						t.Stop()
						return
					}
				case <-rrCtx.Done():
					t.Stop()
					return
				}
			}
		}(l.GeekCard.GeekName, l.GeekCard.SecurityID)

		time.Sleep(5 * time.Second)
	}

	wg.Wait()

	utils.SendFeiShu("====Boss本次招聘任务结束====")
}

// 根据JobId搜索候选人
func (boss *Boss) searchGeekByJobId(jobId, jobName string) ([]*model.Geek, error) {
	var geeks []*model.Geek
	geekList, err := boss.ListRecommend(jobId)
	if err != nil {
		return nil, err
	}
	if len(geekList) == 0 {
		utils.SendFeiShu("Boss当前需要重新验证")
		return nil, accountUnusual
	}

	for _, geek := range geekList {
		logs.Printf("候选人: %s  期待职位：%s \n", geek.GeekCard.GeekName, geek.GeekCard.ExpectPositionName)
		if boss.selectGeek(geek, jobName) {
			logs.Printf("候选人: %s  进入队列, 分值: %d\n", geek.GeekCard.GeekName, geek.Weight)
			geeks = append(geeks, geek)
		}
	}
	return geeks, nil
}

// 筛选并打分
func (boss *Boss) selectGeek(geek *model.Geek, jobName string) bool {
	// 已经打过招呼了
	if geek.HaveChatted == 1 {
		return false
	}
	// 已经被同事撩过
	if geek.Cooperate == constant.CommunicatedYes {
		return false
	}
	// 岗位匹配
	if !utils.MatchJob(jobName, geek) {
		return false
	}
	//  是否是本科
	if geek.GeekCard.GeekDegree == "本科" {
		geek.Weight += 3
	}
	//  是否是硕士
	if geek.GeekCard.GeekDegree == "硕士" {
		geek.Weight += 4
	}
	// 是否是211
	if utils.IsContains(config.School211, geek.GeekCard.GeekEdu.School) {
		geek.Weight += 3
	}
	// 是否是985
	if utils.IsContains(config.School985, geek.GeekCard.GeekEdu.School) {
		geek.Weight += 4
	}
	// 是否在大厂
	for _, w := range geek.GeekCard.GeekWorks {
		if utils.IsContains(config.GoodCompany, w.Company) {
			geek.Weight += 5
			break
		}
	}
	// 工作年限大于3年
	workStr := strings.ReplaceAll(geek.GeekCard.GeekWorkYear, "年", "")
	if years, err := strconv.Atoi(workStr); err == nil && years >= 3 {
		geek.Weight += 3
	}
	// 年龄
	ageStr := strings.ReplaceAll(geek.GeekCard.AgeDesc, "岁", "")
	if age, err := strconv.Atoi(ageStr); err == nil && age >= 26 && age <= 35 {
		geek.Weight += 2
	}
	// 在职-月内到岗
	if strings.Contains(geek.GeekCard.ApplyStatusDesc, "月内到岗") {
		geek.Weight += 2
	}
	// 离职-随时到岗
	if strings.Contains(geek.GeekCard.ApplyStatusDesc, "离职") {
		geek.Weight += 3
	}
	// 今日活跃
	if strings.Contains(geek.ActiveTimeDesc, "今日活跃") {
		geek.Weight += 1
	}
	// 刚刚活跃
	if strings.Contains(geek.ActiveTimeDesc, "刚刚活跃") {
		geek.Weight += 2
	}
	return true
}

// 打招呼
// 需要设置自动打招呼
func (boss *Boss) hello(jobId string, geekId int, encryptGeekId, lid, securityId string, expectId int) error {
	if _, ok := boss.talked.Load(geekId); ok {
		return ErrTalked
	}
	if utils.MapLen(&boss.talked) == boss.Jobs[jobId].HelloNum {
		return maxLimit
	}
	uri := fmt.Sprintf("https://www.zhipin.com/wapi/zpboss/h5/chat/start?_=%d", time.Now().Unix())
	urlQuery := url.Values{}
	urlQuery.Add("jid", jobId)
	urlQuery.Add("gid", encryptGeekId)
	urlQuery.Add("lid", lid)
	urlQuery.Add("expectId", fmt.Sprintf("%d", expectId))
	urlQuery.Add("securityId", securityId)

	data := strings.NewReader(urlQuery.Encode())
	req, _ := http.NewRequest(http.MethodPost, uri, data)
	boss.addHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("hello request", err.Error())
		return err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	str := string(bs)
	if strings.Contains(str, "今日沟通已达上限") {
		return maxLimit
	}
	// 标记已经打过招呼了
	boss.talked.Store(geekId, 1)

	return nil
}

// 接收简历
func (boss *Boss) acceptResumes(mid, securityId string) error {
	uri := "https://www.zhipin.com/wapi/zpchat/exchange/accept"
	urlQuery := url.Values{}
	urlQuery.Add("mid", mid)
	urlQuery.Add("type", constant.RequestTypeToMe)
	urlQuery.Add("securityId", securityId)

	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(urlQuery.Encode()))
	boss.addHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("acceptResumes request", err.Error())
		return err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bs))
	return nil
}

// 向牛人请求简历
// 每隔一段时间请求一次，直到对方回复我们，建立好友关系为止
func (boss *Boss) requestResumes(name, securityId string) error {
	uri := "https://www.zhipin.com/wapi/zpchat/exchange/request"
	urlQuery := url.Values{}
	urlQuery.Add("type", constant.RequestTypeToGeek)
	urlQuery.Add("securityId", securityId)

	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(urlQuery.Encode()))
	boss.addHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(bs), "好友关系校验失败") {
		return notFriend
	}
	var temp *model.RequestResumesResp
	if err = json.Unmarshal(bs, &temp); err != nil {
		return err
	}
	if fmt.Sprintf("%d", temp.ZpData.Type) != constant.RequestTypeToGeek {
		return notFriend
	}
	logs.Printf("请求候选人:%s的简历成功 \n", name)
	return nil
}

// 获取推荐牛人列表
func (boss *Boss) ListRecommend(jobId string) ([]*model.Geek, error) {
	uri := fmt.Sprintf("https://www.zhipin.com/wapi/zprelation/interaction/bossGetGeek?")
	urlQueue := url.Values{}
	urlQueue.Add("gender", "0")
	urlQueue.Add("exchangeResumeWithColleague", "0")
	urlQueue.Add("switchJobFrequency", "0")
	urlQueue.Add("activation", "0")
	urlQueue.Add("recentNotView", "0")
	urlQueue.Add("school", "0")
	urlQueue.Add("major", "0")
	urlQueue.Add("experience", "0")
	urlQueue.Add("jobid", jobId)
	urlQueue.Add("degree", "0")
	urlQueue.Add("salary", "0")
	urlQueue.Add("intention", "0")
	urlQueue.Add("refresh", fmt.Sprintf("%d", time.Now().Unix()))
	urlQueue.Add("status", "1")
	urlQueue.Add("cityCode", "")
	urlQueue.Add("businessId", "0")
	urlQueue.Add("source", "")
	urlQueue.Add("districtCode", "0")
	urlQueue.Add("page", fmt.Sprintf("%d", 1))
	urlQueue.Add("tag", "1")

	uri = uri + urlQueue.Encode()
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	boss.addHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		logs.Println("ListRecommend request", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(bs), "登录状态已失效") {
		return nil, notLogin
	}
	var temp *model.GeekListResp
	err = json.Unmarshal(bs, &temp)
	if err != nil {
		return nil, err
	}
	return temp.ZpData.GeekList, nil
}

func (boss *Boss) addHeader(req *http.Request) {
	req.Header.Add("cookie", boss.Cookie)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
}

// 设置自动打招呼语
// 根据Job设置
func (boss *Boss) SetHelloMsg() {
	// 开启自动打招呼
	uri := "https://www.zhipin.com/wapi/zpchat/greeting/updateGreeting"
	values := url.Values{}
	values.Add("status", "1")
	values.Add("templateId", "")
	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(values.Encode()))
	boss.addHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("open auto greeting", err.Error())
		return
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(bs), "Success") {
		logs.Println("已开启自动打招呼")
	}
	// 获取职位列表
	uri = "https://www.zhipin.com/wapi/zpchat/greeting/job/get"
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	boss.addHeader(req)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("setHelloMsg get", err.Error())
		return
	}
	defer resp.Body.Close()

	bs, _ = ioutil.ReadAll(resp.Body)
	var t *model.JobHelloMsg
	err = json.Unmarshal(bs, &t)
	if err != nil {
		logs.Println("unmarshal job message", err.Error())
		return
	}

	// 设置每个岗位的打招呼语
	uri = "https://www.zhipin.com/wapi/zpchat/greeting/job/save"
	for _, job := range t.ZpData.Jobs {
		// 如果设置过了,就不再设置了
		if job.JobGreeting != "" {
			continue
		}
		data := url.Values{}
		data.Add("encJobId", job.EncJobID)
		data.Add("encGreetingId", job.EncGreetingID)
		data.Add("content", fmt.Sprintf("你好，这边是得物APP，我们目前正在大力扩招%s，如果您有兴趣的话，方便发一份简历给我吗？期待你的加入～", job.JobName))

		req, _ = http.NewRequest(http.MethodPost, uri, strings.NewReader(data.Encode()))
		boss.addHeader(req)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			logs.Println("save job hell msg", err.Error())
			continue
		}
		defer resp.Body.Close()

		bs, _ := ioutil.ReadAll(resp.Body)
		if strings.Contains(string(bs), "Success") {
			logs.Printf("设置职位: %s 的打招呼语成功", job.JobName)
		}
	}
}

//  TODO: 扫码登录,暂不支持
func (boss *Boss) GetQRId(ctx context.Context) {
	// 取qrId
	uri := "https://login.zhipin.com/wapi/zppassport/captcha/randkey"
	values := url.Values{}
	values.Add("pk", "cpc_user_sign_up")
	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(values.Encode()))
	boss.addHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("get qr id", err.Error())
		return
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	var msg *model.QRMsg
	if err = json.Unmarshal(bs, &msg); err != nil {
		logs.Println("unmarshal qr msg", err.Error())
		return
	}
	// 取qrId
	qrId := msg.ZpData.QrID

	newCtx, _ := context.WithTimeout(ctx, 10*time.Minute)
	go func(qrId string) {
		t := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-t.C:
				// 获取set-cookie
				if err := boss.setCookie(qrId); err == nil {
					t.Stop()
					return
				}
			case <-newCtx.Done():
				return
			}
		}

	}(qrId)

}

func (boss *Boss) setCookie(qrId string) error {
	uri := "https://login.zhipin.com/wapi/zppassport/qrcode/dispatcher?"
	values := url.Values{}
	values.Add("qrId", qrId)
	values.Add("_", fmt.Sprintf("%d", time.Now().Unix()))
	req, _ := http.NewRequest(http.MethodGet, uri, strings.NewReader(values.Encode()))
	boss.addHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("get cookie", err.Error())
		return err
	}
	defer resp.Body.Close()

	var setCookie string
	for _, c := range resp.Header["Set-Cookie"] {
		setCookie += c
		setCookie += ";"
	}
	if setCookie == "" {
		return fmt.Errorf("no cookie")
	}
	return nil
}

// 获取job列表
func (boss *Boss) ListJobs() []*model.Job {
	uri := "https://www.zhipin.com/wapi/zpjob/job/data/list?"
	values := url.Values{}
	values.Add("position", "0")
	values.Add("searchStr", "0")
	values.Add("page", "1")
	values.Add("_", fmt.Sprintf("%d", time.Now().Unix()))

	req, _ := http.NewRequest(http.MethodGet, uri, strings.NewReader(values.Encode()))
	boss.addHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Println("list job", err.Error())
		return nil
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)

	var jResp *model.JobListResp
	if err = json.Unmarshal(bs, &jResp); err != nil {
		logs.Println("unmarshal list job", err.Error())
		return nil
	}
	var jobs []*model.Job
	for _, j := range jResp.ZpData.Data {
		if j.JobStatus == 0 {
			jobs = append(jobs, &model.Job{
				JobId:   j.EncryptJobID,
				JobName: j.JobName,
			})
		}
	}
	return jobs
}
