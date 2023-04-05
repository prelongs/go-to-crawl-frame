package gofastdfs

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/eventials/go-tus"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"os"
)

func Upload(path string) {
	logger := g.Log()

	var obj interface{}
	req := httplib.Post("http://15.235.66.41:8080/group1/upload")
	req.PostFile("file", path)         // 本地需要上传的文件的完整路径名
	req.Param("output", "json")        // API调用后响应体格式
	req.Param("filename", "test3.txt") // 给文件重新命名
	req.Param("scene", "test")         // 场景：会转化为go-fastdfs的文件根目录
	req.Param("path", "a/b/c")         // 自定义文件保存到go-fastdfs的路径，设置后会忽略scene设置
	// 参数综合建议：
	// 1、指定file，output，path即可，path可根据世纪业务场景做迁移工作，如果利用go-fastdfs默认文件上传后的目录组织方式(日期方式)，不便于管理，除非
	//    一开始设计之初就使用日期的组织方式
	req.ToJSON(&obj)
	logger.Info(gctx.GetInitCtx(), obj)
}

// UploadBreakpoint 断点续传
func UploadBreakpoint(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// create the tus client.
	client, err := tus.NewClient("http://15.235.66.41:8080/big/upload", nil)
	fmt.Println(err)
	// create an upload from a file.
	upload, err := tus.NewUploadFromFile(f)
	fmt.Println(err)
	// create the uploader.
	uploader, err := client.CreateUpload(upload)
	fmt.Println(err)
	// start the uploading process.
	fmt.Println(uploader.Upload())
}
