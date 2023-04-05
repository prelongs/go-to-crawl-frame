package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
)

//跨域处理中间件
func CORS(r *ghttp.Request) {
	fmt.Println("跨域处理中间件")
	r.Response.CORSDefault()
	// 前置中间件
	r.Middleware.Next()
}
