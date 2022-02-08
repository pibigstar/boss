package main

import "github.com/pibigstar/boss/online"

func main() {

	online.RunHttp(8080)

	// 1. 读取可用用户

	// 2. 读取用户配置的分值计算规则

	// 3. 读取211学校、985学校、大厂信息

	// 4. 读取用户配置的需要招聘的job
	// 每个job抓取几分钟，对几个人进行打招呼
	// 最大的候选人队列值

	// 5. for 对 job 进行招聘
}
