package boot

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/swagger/v2"
)

// 用于应用初始化。
func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})

	gview.Instance().BindFuncMap(gview.FuncMap{})
}
