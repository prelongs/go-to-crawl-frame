package middleware

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/app/dao"
	"github.com/JervisPG/go-to-crawl-frame/app/model"
	"github.com/JervisPG/go-to-crawl-frame/app/utils"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

func OperLog(r *ghttp.Request) {
	// 后置中间件
	r.Middleware.Next()
	// 中间件处理逻辑
	fmt.Println("操作日志中间件")

	// 分析请求URL地址
	urlArr := gstr.Split(r.URL.String(), "?")
	urlItem := gstr.Split(urlArr[0], "/")
	if len(urlItem) < 4 {
		return
	}
	operType := urlItem[3]
	// 拼接节点
	permission := "sys:" + urlItem[2] + ":" + operType
	// 查询节点信息
	info, err := dao.Menu.Where("permission=?", permission).FindOne()
	if err != nil || info == nil {

		return
	}

	// 创建日志对象
	var entity model.OperLog
	entity.Model = info.Title
	if operType == "add" || operType == "addz" {
		// 新增
		entity.OperType = 1
	} else if operType == "update" {
		// 修改
		entity.OperType = 2
	} else if operType == "delete" || operType == "dall" {
		// 删除
		entity.OperType = 3
	} else if operType == "list" {
		// 查询
		entity.OperType = 4
	} else if operType == "status" {
		// 设置状态
		entity.OperType = 5
	} else if operType == "import" {
		// 导入
		entity.OperType = 6
	} else if operType == "export" {
		// 导出
		entity.OperType = 7
	} else if operType == "permission" {
		// 设置权限
		entity.OperType = 8
	} else if operType == "resetPwd" {
		// 设置密码
		entity.OperType = 9
	} else {
		// 其他
		entity.OperType = 0
	}

	entity.OperMethod = r.Method
	entity.OperName = utils.UInfo(r).Realname
	entity.Username = utils.UInfo(r).Username
	entity.OperUrl = r.URL.String()
	entity.OperIp = r.GetClientIp()
	entity.OperLocation = utils.GetIpCity(entity.OperIp)
	entity.RequestParam = string(r.GetBody())
	entity.Status = 0
	entity.UserAgent = r.UserAgent()
	entity.CreateUser = utils.Uid(r)
	entity.CreateTime = gtime.Now()
	entity.Mark = 1
	dao.OperLog.Insert(entity)
}
