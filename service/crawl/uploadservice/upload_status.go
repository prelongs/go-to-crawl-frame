package uploadservice

const (
	//上传状态.0-创建任务;1-上传中;2-上传完成;3-流媒体处理中;4-流媒体处理结束;5-流媒体处理异常;6-成功請求CMS
	Init           = 0
	Uploading      = 1
	Uploaded       = 2
	Transforming   = 3
	Transformed    = 4
	TransformErr   = 5
	CmsPostSuccess = 6
)

const (
	FileTypeVideo   = 1
	FileTypeBigFile = 2
)
