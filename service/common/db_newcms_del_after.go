package common

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

//删除表格后的回调命名方法Del_tablename
//删除指定表成功后会执行对应方法 用于后续缓存的更新或一对多的处理
//参数为删除的ID
func (cms *sDbNewCms) Del_capp(delId interface{}) {
	var tableName = cms.prefix + "cappinfo"
	_, _ = cms.db.Model(tableName).Where(g.Map{"appid": delId}).Delete()
	fmt.Printf("------------------%v-----------------\r\n", delId)
}
