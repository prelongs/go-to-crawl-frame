package lockservice

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-frame/service/configservice"
	"time"
)

const (
	CrawlSeleniumLock = "crawl:selenium:lock"
)

// 不管哪个需要selenium的任务，只有获取到锁才能继续往下走，避免资源竞争，后期可区分端口做并发使用selenium
func TryLockSelenium() bool {
	return tryLockSeleniumMinutes(1)
}

// 同TryLockSelenium，selenium耗时比较长的任务锁3分钟
func TryLockSeleniumLong() bool {
	return tryLockSeleniumMinutes(3)
}

func tryLockSeleniumMinutes(minutes int) bool {
	value, _ := gcache.SetIfNotExist(gctx.GetInitCtx(), CrawlSeleniumLock, new(struct{}), gtime.M*time.Duration(minutes))
	return value
}

func AddSeleniumExpireTime(seconds int) {
	expireTime, _ := gcache.GetExpire(gctx.GetInitCtx(), CrawlSeleniumLock)
	if expireTime == -1 {
		// 已经过期的无需延长
		return
	}
	expireTime = expireTime + time.Second*time.Duration(seconds)
	gcache.UpdateExpire(gctx.GetInitCtx(), CrawlSeleniumLock, expireTime)
}

// 释放锁，等待别的定时任务下次执行可抢占到锁
func ReleaseLockSelenium() {
	_, _ = gcache.Remove(gctx.GetInitCtx(), CrawlSeleniumLock)
}

func IncreaseValue(key string) bool {
	cacheValue, _ := gcache.Get(gctx.GetInitCtx(), key)
	if cacheValue.Val() == nil {
		gcache.Set(gctx.GetInitCtx(), key, 1, time.Duration(0))
	} else {
		max := configservice.GetInt(key, 10)
		increasedValue := cacheValue.Int() + 1
		if increasedValue > max {
			return false
		} else {
			gcache.Set(gctx.GetInitCtx(), key, increasedValue, time.Duration(0))
		}
	}

	return true
}

func DecreaseValue(key string) bool {
	cacheValue, _ := gcache.Get(gctx.GetInitCtx(), key)
	if cacheValue.Val() == nil {
		return false
	} else {
		current := cacheValue.Int()
		if current <= 0 {
			return false
		} else {
			gcache.Set(gctx.GetInitCtx(), key, current-1, time.Duration(0))
			return true
		}
	}

}
