package model

type ScoreConfig struct {
	Score211     int `json:"score211"`     // 211
	Score958     int `json:"score985"`     // 985
	GoodCompany  int `json:"goodCompany"`  // 大公司
	Undergrad    int `json:"undergrad"`    // 本科
	Master       int `json:"master"`       // 硕士
	WorkTime     int `json:"workTime"`     // 工作大于3年
	AgeOver35    int `json:"ageOver35"`    // 年龄大于35(可设置负值)
	OnlineWork   int `json:"onlineWork"`   // 在职
	OfflineWork  int `json:"offlineWork"`  // 离职
	ActiveToday  int `json:"activeToday"`  // 今日活跃
	ActiveMinute int `json:"activeMinute"` // 刚刚活跃
}
