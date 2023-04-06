package sysservice

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/dao"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	NsProxy       = "crawlerProxy"
	NsCountryCode = "countryCode"
)

const (
	DictEnable  = 1
	DictDisAble = 0
)

var (
	cdc = dao.CmsCrawlDict.Columns()
)

func GetRandomProxyUrl() string {
	var one *entity.CmsCrawlDict
	_ = dao.CmsCrawlDict.Ctx(gctx.GetInitCtx()).
		Where(cdc.Namespace, NsProxy).
		Where(cdc.DictStatus, DictEnable).
		OrderRandom().Scan(&one)

	if one == nil {
		return ""
	}

	return fmt.Sprintf("http://%s", one.DictValue)
}
