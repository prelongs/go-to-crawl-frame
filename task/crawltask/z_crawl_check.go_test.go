package crawltask

import (
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestCrawlCheckTask(t *testing.T) {
	CrawlCheckTask.CheckQQLoginTask(gctx.GetInitCtx())
}
