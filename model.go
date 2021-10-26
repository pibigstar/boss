package main

type User struct {
	Id         int          `json:"id"`
	UserName   string       `json:"username"`
	Cookie     string       `json:"cookie"`
	Status     int          `json:"status"`
	ExtraInfos []*ExtraInfo `json:"extraInfos"`
	Jobs       []*JobDetail `json:"jobs"`
}

type JobDetail struct {
	JobId     int    `json:"jobId"`
	Name      string `json:"name"`      // 职位名
	ShortName string `json:"shortName"` // 职位简写
	GreetNum  int    `json:"greetNum"`  // 打招呼人数
}

type ExtraInfo struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}
