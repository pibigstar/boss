package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pibigstar/boss/logs"
)

type Config struct {
	Online bool   `json:"online"`
	Mysql  *Mysql `json:"mysql"`
}

type Mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

var (
	cfg        *Config
	defaultCfg = "config.json"
)

func GetConfig() *Config {
	return cfg
}

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(currentFile)
	bs, err := os.ReadFile(filepath.Join(basePath, defaultCfg))
	if err != nil {
		logs.Println("failed to read config.json", err.Error())
		return
	}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		if err != nil {
			logs.Println("failed to unmarshal config.json", err.Error())
			return
		}
	}
	// 在线环境的话，需要连接mysql
	if cfg.Online && cfg.Mysql != nil {
		connectMysql()
	}
}

// ====== db相关操作 ======
var db *sql.DB

//  获取mysql连接
func connectMysql() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4,utf8",
		cfg.Mysql.UserName,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DB)

	var err error
	db, err = sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}
