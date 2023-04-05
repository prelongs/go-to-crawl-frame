package crawldto

import (
	"go-to-crawl-frame/db/mysql/model/entity"
)

type CmsCrawl struct {
	*entity.CmsCrawlQueue
	ShowStatus   int    `json:"showStatus"`
	ResourcePath string `json:"resourcePath"`
}
