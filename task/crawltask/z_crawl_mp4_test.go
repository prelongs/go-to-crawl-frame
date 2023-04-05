package crawltask

import (
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestDownloadMp4Task(t *testing.T) {
	DownloadMp4Task(gctx.GetInitCtx())
}

func TestDownloadMp4Type1Task(t *testing.T) {
	DownloadMp4Type1Task(gctx.GetInitCtx())
}
