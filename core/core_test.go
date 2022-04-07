package core

import (
	"context"
	"testing"

	"github.com/pibigstar/boss/model"
)

var boss = &Boss{
	User: model.User{
		UserName: "测试",
	},
	Jobs: map[string]model.Job{
		"3b7ba39cd5535e7a1n1_09u4FVVT": DefaultJob("3b7ba39cd5535e7a1n1_09u4FVVT", "Golang"),
	},
	Cookie: `lastCity=101020100; wd_guid=adaabd69-ccd3-4d22-baad-668677ef2d08; historyState=state; __g=-; Hm_lvt_194df3105ad7148dcf2b98a91b5e727a=1649315746; Hm_lpvt_194df3105ad7148dcf2b98a91b5e727a=1649315755; wt2=DUwxBaTxqlL90CsvpBBCc2N8cbEEgQxeI5QvQT0JiFo--YbRwsI3McHplaItKUTMlyXR7mEPe0jfOdsCnHPOTBQ~~; wbg=1; acw_tc=0bdd34c216493211266456865e01aa4f03b9fdd5cff26ed36e19f7c43b704d; zp_token=V1QtMhFuf52VxgXdNtxh4aKyKx7DPfwg~~; __l=l=/www.zhipin.com/web/boss/job/list&s=3&friend_source=0&s=3&friend_source=0; _dd_s=logs=1&id=4a9715ef-473a-4f88-9517-198f879fd962&created=1649322326781&expire=1649323271898; __c=1649315744; __a=47595245.1586530028.1636943029.1649315744.497.5.11.488`,
}

func TestMain(m *testing.M) {
	m.Run()
}

// 测试招人
func TestHiring(t *testing.T) {
	boss.Hiring()
}

func TestListRecommend(t *testing.T) {
	for jobId := range boss.Jobs {
		geeks, err := boss.listRecommend(jobId)
		if err != nil {
			t.Error(err)
		}
		for _, geek := range geeks {
			t.Log(geek.GeekCard.GeekName)
		}
	}
}

func TestSetHelloMsg(t *testing.T) {
	boss.SetHelloMsg()
}

func TestGetQR(t *testing.T) {
	boss.getQRId(context.Background())
}

func TestListJob(t *testing.T) {
	jobs := boss.ListJobs()
	for _, j := range jobs {
		t.Log(j.JobName, j.JobId)
	}
}
