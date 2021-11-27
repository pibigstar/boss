package offline

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pibigstar/boss/core"
	"github.com/pibigstar/boss/logs"
)

var boss = &core.Boss{
	Jobs: make(map[string]core.Job),
}

func RunBoss() {
	initBoss()
	if len(boss.Jobs) == 0 {
		inputJobs()
	}
	if len(boss.Jobs) == 0 {
		fmt.Println("暂时没有需要沟通的职位~")
		return
	}
	boss.Hiring()
}

// 输入存储job信息
func inputJobs() {
	jobs := boss.ListJobs()
	if len(jobs) == 0 {
		logs.Println("你没有开放的职位")
		return
	}
	for i, j := range jobs {
		fmt.Printf("编号:%d 职位: %s \n", i, j.JobName)
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入你要沟通的职位编号:")

	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println("输入错误!")
		return
	}
	ids := strings.Split(input, ",")
	var str string
	for _, id := range ids {
		id = strings.ReplaceAll(id, "\n", "")
		id = strings.TrimSpace(id)
		i, err := strconv.Atoi(id)
		if err != nil || i >= len(jobs) {
			fmt.Printf("%s 编号有误 \n", id)
			continue
		}
		jobId := jobs[i].JobId
		jobName := jobs[i].JobName
		boss.Jobs[jobId] = core.DefaultJob(jobId, jobName)
		str += fmt.Sprintf("%s   //%s \n", jobId, jobName)
	}

	// 储存
	err = ioutil.WriteFile(jobsFile, []byte(str), 0666)
	if err != nil {
		logs.Println("存储jobId信息失败", err.Error())
	}
}

// 监听cookie变化
func watchCookie() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	err = watcher.Add(cookieFile)
	if err != nil {
		logs.Println("watch cookie.txt", err.Error())
		return
	}
	// 开始监听
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				// cookie文件有变动，重新设置cookie
				readCookie()

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logs.Println("watcher error:", err)
			}
		}
	}()
}
