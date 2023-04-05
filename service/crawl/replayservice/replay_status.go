package replayservice

const (
	// 抓取状态.0-创建任务;1-抓取中;2-抓取失败;3-抓取完成;4-发布cms完成
	ManifestTaskInit             = 0
	ManifestTaskCrawling         = 1
	ManifestTaskCrawlErr         = 2
	ManifestTaskCrawlFinish      = 3
	ManifestTaskCrawlPostSuccess = 4

	//20-转码中
	ManifestTaskCrawlTraning = 20
)

const (
	// 抓取状态.0-创建任务;1-录制中;2-录制失败;3-录制完成;4-切片中;5-切片失败;6-切片成功
	ProgramTaskInit        = 0
	ProgramTaskCrawling    = 1
	ProgramTaskCrawlErr    = 2
	ProgramTaskCrawlFinish = 3
	ProgramTaskParsing     = 4
	ProgramTaskParseErr    = 5
	ProgramTaskParsed      = 6
)

//`{"code":0,"data":{"code":200,"data":"{\"uuid\": \"487895cae7a9c9110f68c90a82d43a8f\"}","error":""}}`
//回看发布返回
type ReplayRes struct {
	Code int         `json:"code" dc:"返回代码0 为成功 其他失败"`
	Data *ReplayData `json:"data" dc:"返回详细信息"`
}
type ReplayData struct {
	Code  int    `json:"code" dc:"200成功获取节目ID 其他失败"`
	Data  string `json:"data" dc:"id详情"`
	Error string `json:"error" dc:"获取ID错误信息"`
}
type ReplayUUID struct {
	UUID string `json:"uuid" dc:"节目uuid"`
}
