package main

import (
	"github.com/pibigstar/boss/config"
	"github.com/pibigstar/boss/offline"
	"github.com/pibigstar/boss/online"
)

func main() {
	if config.GetConfig().Online {
		online.CronRun()
		online.RunHttp(8080)
	} else {
		offline.RunBoss()
	}
}
