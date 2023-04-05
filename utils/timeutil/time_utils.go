package timeutil

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
)

const (
	YYYY_MM_DD            = "2006-01-02"
	YYYY_MM_DD_JOIN       = "20060102"
	YYYY_MM_DD_SLASH      = "2006/01/02"
	YYYY_MM_DD_HH_MM_SS   = "2006-01-02 15:04:05"
	YYYY_MM_DD_HH_MM_JOIN = "200601021504"
	HH_MM_SS              = "15:04:05"
)

func Sleep(seconds int) {
	du := gtime.S * time.Duration(seconds)
	g.Log().Info(gctx.GetInitCtx(), du, "进入睡眠...")
	time.Sleep(du)
	g.Log().Info(gctx.GetInitCtx(), "唤醒睡眠...")
}
