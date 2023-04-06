package crawltask

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/gogf/gf/v2/os/gctx"
	"os/exec"

	"testing"
)

func TestIqiyi(t *testing.T) {
	// GET参数部分可不用录入
	doStartTmpSeed("https://www.iqiyi.com/v_19rra7yxcs.html")
}

func TestBanan(t *testing.T) {
	doStartTmpSeed("https://banan.tv/vodplay/37352-1-1.html")
}

func TestNiVod(t *testing.T) {
	doStartTmpSeed("https://www.nivod.tv/Ok2b8bCuT005mcWwpOKW2ak5AvQZOOQY-afn1j1MYA9UGQBqRdu3YdqW8r9mLO2Un-1080-0-0-play.html?x=1")
}

func TestNiVod3(t *testing.T) {
	doStartTmpSeed("https://www.nivod3.tv/1AyYWd1WFag2bKjJliUuuAQFT2vgDKzB-0-0-0-0-play.html?x=1")
}

func TestQQ(t *testing.T) {
	doStartTmpSeed("https://v.qq.com/x/cover/2w2legt0g8z26al/i0029aelh3d.html")
}

func TestTangRenJie(t *testing.T) {
	doStartTmpSeed("https://www.tangrenjie.tv/vod/play/id/197503/sid/1/nid/31.html")
}

func TestTangRenJie2(t *testing.T) {
	doStartTmpSeed("https://www.tangrenjie.tv/vod/play/id/218662/sid/1/nid/8.html")
}

func TestOleVod(t *testing.T) {
	doStartTmpSeed("https://www.olevod.com/index.php/vod/play/id/33904/sid/1/nid/45.html")
}

func TestNunuyy(t *testing.T) {
	doStartTmpSeed("https://www.nunuyy3.org/dongman/102123.html")
}

func TestNunuyyWithParams(t *testing.T) {
	seed := new(entity.CmsCrawlQueue)
	seed.CrawlSeedUrl = "https://www.nunuyy3.org/dongman/102123.html"
	// 努努大部分资源都支持量子资源，且量子资源的m3u8未加密，因此videoItem需要切换到量子资源下的节目单来摘录到数据库
	seed.CrawlSeedParams = `{"videoItem":"第25集"}`
	DoStartCrawlVodFlow(seed)
}

// 测试通过，按需求对接数据库改成多线程下载就行
func TestBilibili(t *testing.T) {
	doStartTmpSeed("https://www.bilibili.com/video/BV1yg411T7Za?p=7&vd_source=e8bcede57b979b0eed49d9041d869a8e")
}

func TestBilibiliVIP(t *testing.T) {
	doStartTmpSeed("https://www.bilibili.com/bangumi/play/ep672789?spm_id_from=333.6.0.0")
}

func TestCrawlHostTypeNormal(t *testing.T) {
	CrawlTask.CrawlUrlTask(gctx.GetInitCtx())
}

func TestCrawlUrlType1Task(t *testing.T) {
	CrawlTask.CrawlUrlType1Task(gctx.GetInitCtx())
}

func TestCrawlUrlType2Task(t *testing.T) {
	CrawlTask.CrawlUrlType2Task(gctx.GetInitCtx())
}

func doStartTmpSeed(url string) {
	seed := new(entity.CmsCrawlQueue)
	seed.CrawlSeedUrl = url
	DoStartCrawlVodFlow(seed)
}

func TestName(t *testing.T) {
	cmd := exec.Command("curl", "https://www.baidu.com")

	out, err := cmd.Output()
	fmt.Println(out)
	fmt.Println(err)

	c := "curl https://www.baidu.com"

	cmd = exec.Command("sh", "-c", c)

	out, err = cmd.Output()
	fmt.Println(out)
	fmt.Println(err)
}
