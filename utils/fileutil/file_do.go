package fileutil

const (
	POST = "POST"
	GET  = "GET"
)

type DownloadDO struct {
	method     string
	url        string
	proxy      string // 代理地址
	saveFile   string // 保存文件
	retry      int    // 重试次数
	timeout    int64  // 超时
	headers    map[string][]string
	body       string // POST请求体
	preRspBody string // 预响应体
}

type DownloadBuilder struct {
	DownloadAttr *DownloadDO
}

func CreateBuilder() *DownloadBuilder {
	builder := new(DownloadBuilder)
	builder.DownloadAttr = new(DownloadDO)
	builder.Method(GET)
	return builder
}

func (builder *DownloadBuilder) Method(method string) *DownloadBuilder {
	builder.DownloadAttr.method = method
	return builder
}

func (builder *DownloadBuilder) Url(url string) *DownloadBuilder {
	builder.DownloadAttr.url = url
	return builder
}

func (builder *DownloadBuilder) Proxy(proxy string) *DownloadBuilder {
	builder.DownloadAttr.proxy = proxy
	return builder
}

func (builder *DownloadBuilder) SaveFile(saveFile string) *DownloadBuilder {
	builder.DownloadAttr.saveFile = saveFile
	return builder
}

func (builder *DownloadBuilder) Retry(retry int) *DownloadBuilder {
	builder.DownloadAttr.retry = retry
	return builder
}

func (builder *DownloadBuilder) Timeout(timeout int64) *DownloadBuilder {
	builder.DownloadAttr.timeout = timeout
	return builder
}

func (builder *DownloadBuilder) Header(key, value string) *DownloadBuilder {
	if builder.DownloadAttr.headers == nil {
		builder.DownloadAttr.headers = make(map[string][]string, 10)
	}
	values := builder.DownloadAttr.headers[key]
	if values == nil {
		values = []string{}
	}
	builder.DownloadAttr.headers[key] = append(values, value)
	return builder
}

func (builder *DownloadBuilder) Headers(headers map[string][]string) *DownloadBuilder {
	builder.DownloadAttr.headers = headers
	return builder
}

func (builder *DownloadBuilder) Body(body string) *DownloadBuilder {
	builder.DownloadAttr.body = body
	return builder
}

func (builder *DownloadBuilder) PreRspBody(preRspBody string) *DownloadBuilder {
	builder.DownloadAttr.preRspBody = preRspBody
	return builder
}

func (builder *DownloadBuilder) Build() *DownloadDO {
	return builder.DownloadAttr
}
