package ffmpegutil

import "net/http"

type M3u8DO struct {
	Schema          string
	Host            string
	FromUrl         string
	FromUrlProxy    string // FromUrl用数据库配置的Proxy去加载
	FromBaseUrl     string
	FromFile        string
	FromDir         string
	StreamLineList  []StreamLineDO
	SrcLineCount    int
	PlaySecondTotal float32 // 总的播放秒数，可以为浮点型
	PngHeaderSize   int64
	MP4SaveFile     string // 下载MP4后临时变量

	MapSaveFile   string // 统一TS文件头下载路径临时变量
	MapOriginName string // URI="xxx"中的xxx

	KeySaveFile   string // 解密文件下载路径临时变量
	KeyOriginName string // URI="xxx"中的xxx

	Headers *http.Header // TS资源需要header才能下载的场景

	NotDownloadAll bool
}

type StreamLineDO struct {
	LineType           int     // 行类型
	SrcType            int     // TS资源类型
	OriginLine         string  // 原始行
	OriginTsSrcName    string  // 原始TS资源名称
	TransformedLine    string  // 转变后的行
	TransformedSrcName string  // 转变后资源名
	PlaySecond         float32 // 播放秒数，可以为浮点型
}
