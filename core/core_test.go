package core

import (
	"context"
	"testing"

	"github.com/pibigstar/boss/model"
)

var boss = &Boss{
	Jobs: map[string]model.Job{
		"bbdc011c6eb717ad33Z83tq6FlE~": DefaultJob("bbdc011c6eb717ad33Z83tq6FlE~", "Golang"),
	},
	Cookie: `_bl_uid=LFkzUq2jsv2wg9ob25C4bvhd9mOj; lastCity=101020100; wd_guid=adaabd69-ccd3-4d22-baad-668677ef2d08; historyState=state; __g=-; Hm_lvt_194df3105ad7148dcf2b98a91b5e727a=1636854418,1636943029; __f=be66bc1f7784b2763f280eb60bb76060; wt2=DXTJ7-i6IHj_J9FsrRjYjF1KIVbYKRfN1XkGURumT1KOSeGBR5snQEJ1Rlc22WgkzYluXIbjfipXuX9wml-7CrA~~; geek_zp_token=V1QtMhFuf52VxgXdNqyBgcICu46DrXwg~~; Hm_lpvt_194df3105ad7148dcf2b98a91b5e727a=1637548544; __l=l=/www.zhipin.com/web/boss/index&r=&g=&s=3&friend_source=0&s=3&friend_source=0; acw_tc=0bdccfc616379059065057791e135eeae0fa2fd6b2bb6fb798502fc03adb60; zp_token=V1QtMhFuf52VxgXdNqyBQYLi217T_Qxw~~; __c=1636943029; __a=47595245.1586530028.1591944613.1636943029.479.4.26.470`,
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
