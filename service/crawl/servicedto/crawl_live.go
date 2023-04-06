package servicedto

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

const (
	defaultWaitSeconds = 3
)

type CrawlLiveDO struct {
	*entity.CmsCrawlLiveConfig
	NextLoadTime *gtime.Time
}

func (r *CrawlLiveDO) LoadNextTimeDefault() {
	r.NextLoadTime = gtime.Now().Add(time.Second * time.Duration(defaultWaitSeconds))
}

func (r *CrawlLiveDO) LoadNextTimeAfterSeconds(floatSeconds float32) {
	// 提前一些时间，一定程度防止丢帧
	r.NextLoadTime = gtime.Now().Add(time.Second * time.Duration(gconv.Int(floatSeconds*2/3)))
}
