package offline

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pibigstar/boss/constant"

	"github.com/pibigstar/boss/core"
)

var (
	cookieFile    = "cookie.txt"
	school985File = "985.txt"
	school211File = "211.txt"
	companyFile   = "company.txt"
	bossLog       = "boss.log"
	jobsFile      = "jobs.txt"
)

func initConfig() {
	// 设置当前运行目录
	setFilePath()
	// 读取cookie信息
	readCookie()
	// 监听cookie文件
	watchCookie()
	// 读取jobId
	readJobs()
	// 读取学校信息
	readSchool()
	// 读取大厂信息
	readCompany()
	// 设置自动打招呼语
	boss.SetHelloMsg()
}

func setFilePath() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(currentFile)

	cookieFile = filepath.Join(basePath, cookieFile)
	school985File = filepath.Join(basePath, school985File)
	school211File = filepath.Join(basePath, school211File)
	jobsFile = filepath.Join(basePath, jobsFile)
	companyFile = filepath.Join(basePath, companyFile)
	//logFile, _ = os.OpenFile(filepath.Join(basePath, bossLog), os.O_RDWR|os.O_CREATE, 0664)
}

func readCookie() {
	bs, _ := ioutil.ReadFile(cookieFile)
	boss.Cookie = string(bs)
}

func readSchool() {
	bs, _ := ioutil.ReadFile(school985File)
	br := bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		boss.ExtraInfo[constant.School985] = append(boss.ExtraInfo[constant.School985], string(a))
	}

	bs, _ = ioutil.ReadFile(school211File)
	br = bufio.NewReader(bytes.NewReader(bs))
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		boss.ExtraInfo[constant.School211] = append(boss.ExtraInfo[constant.School211], string(a))
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
		boss.ExtraInfo[constant.GoodCompany] = append(boss.ExtraInfo[constant.GoodCompany], string(a))
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
			boss.Jobs[jobId] = core.DefaultJob(jobId, jobName)
		}
	}
}
