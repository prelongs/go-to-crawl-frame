package cmsdto

// 分页查询条件
type ListReq struct {
	TableName string `p:"TableName"` // 表名
	Page      int    `p:"page"`      // 页码
	Limit     int    `p:"limit"`     // 每页数
	Where     string `p:"where"`     //查询条件 json字符串 内部查询转换g.map
	Order     string `p:"order"`     //排序条件
}

type UpdateReq struct {
	TableName string `p:"TableName"` // 表名
	Where     string `p:"where"`     //更新条件 json字符串 内部转换g.map
	Data      string `p:"data"`      //更新内容 json字符串 内部转换为g.map
	UserId    int    `p:"userid"`    //更新的用户ID
}

type AddReq struct {
	TableName string `p:"TableName"` // 表名
	Data      string `p:"data"`      //新增 json字符串 内部转换为g.map
	UserId    int    `p:"userid"`    //更新的用户ID
}
type DelReq struct {
	TableName string `p:"TableName"` // 表名
	Where     string `p:"where"`     //删除条件 json字符串 内部转换为g.map
}
