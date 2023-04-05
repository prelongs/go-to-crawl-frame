package httputil

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"go-to-crawl-frame/entity/cmsdto"
	"go-to-crawl-frame/utils/commonutil"
)

func ParseParam(req *ghttp.Request, dto interface{}) {
	err := req.Parse(dto)
	if err != nil {
		Error(req, err.Error())
	}
}

func ParsePageParam(r *ghttp.Request, dto interface{}) {
	ParseParam(r, dto)
	parser, ok := dto.(cmsdto.PageParam)
	if ok {
		parser.InitPageParam()
	}
}

func Error(r *ghttp.Request, msg string) {
	r.Response.WriteJsonExit(commonutil.JsonResult{
		Code: -1,
		Msg:  msg,
	})
}

func Success(r *ghttp.Request) {
	SuccessData(r, nil)
}

func SuccessData(r *ghttp.Request, Data interface{}) {
	SuccessMsgData(r, "SUCCESS", Data)
}

func SuccessMsgData(r *ghttp.Request, msg string, Data interface{}) {
	r.Response.WriteJsonExit(commonutil.JsonResult{
		Code: 0,
		Msg:  msg,
		Data: Data,
	})
}
