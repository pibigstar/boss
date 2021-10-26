package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mysql *Mysql `json:"mysql"`
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

func init() {
	bs, err := os.ReadFile(defaultCfg)
	if err != nil {
		log.Println("failed to read config.json", err.Error())
		return
	}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		if err != nil {
			log.Println("failed to unmarshal config.json", err.Error())
			return
		}
	}
	if cfg.Mysql != nil {
		connectMysql()
	}
}
