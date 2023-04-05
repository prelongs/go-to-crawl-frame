package sysservice

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
	"go-to-crawl-frame/service/crawl/vodservice"
	netUrl "net/url"
	"strings"
)

var (
	C      = dao.CmsCrawlProxy.Columns()
	regTop = "[^.]+\\.(com.cn|com|net.cn|net|org.cn|org|gov.cn|gov|cn|mobi|me|info|name|biz|cc|tv|asia|hk|网络|公司|中国)"
)

func GetProxyByUrl(requestUrl string) string {

	var one *entity.CmsCrawlProxy
	hostIp := configservice.GetString(fmt.Sprintf("crawl.%s", "hostIp"))
	_ = dao.CmsCrawlProxy.Ctx(gctx.GetInitCtx()).
		Where(C.TargetDomain, hostIp).
		Where(C.ProxyStatus, vodservice.CrawProxyOpen).Scan(&one)

	if one == nil {
		return ""
	}

	return one.ProxyUrl
}
func GetProxyByUrlBak(requestUrl string) string {

	if requestUrl == "" {
		return ""
	}

	url, err := netUrl.Parse(requestUrl)
	if err != nil {
		return ""
	}

	host := url.Host
	index := strings.LastIndex(host, ":")
	if index > 0 {
		host = gstr.SubStr(host, 0, index)
	}
	matches, _ := gregex.MatchString(regTop, host)
	var one *entity.CmsCrawlProxy
	_ = dao.CmsCrawlProxy.Ctx(gctx.GetInitCtx()).
		Where(C.TargetDomain, matches[0]).
		Where(C.ProxyStatus, vodservice.CrawProxyOpen).
		Scan(&one)
	if one != nil {
		return one.ProxyUrl
	}
	return ""
}
