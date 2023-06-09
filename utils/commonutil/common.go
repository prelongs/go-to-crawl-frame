package commonutil

type BunissType int

//业务类型
const (
	BOther BunissType = 0 //0其它
	BAdd   BunissType = 1 //1新增
	BEdit  BunissType = 2 //2修改
	BDel   BunissType = 3 //3删除
)

// 返回结果对象
type JsonResult struct {
	Code  int         `json:"code"`   // 响应编码：0成功 401请登录 403无权限 500错误
	Msg   string      `json:"msg"`    // 消息提示语
	AddID int64       `json:"add_id"` //新增成功时最后一条记录的ID
	Data  interface{} `json:"data"`   // 数据对象
	Count int         `json:"count"`  // 记录总数
	Btype BunissType  `json:"btype"`  // 业务类型
}

// 验证码
type CaptchaRes struct {
	Code  int         `json:"code"`  //响应编码 0 成功 500 错误 403 无权限
	Msg   string      `json:"msg"`   //消息
	Data  interface{} `json:"data"`  //数据内容
	IdKey string      `json:"idkey"` //验证码ID
}

// 部门类型
var DEPT_TYPE_LIST = map[int]string{
	1: "公司",
	2: "子公司",
	3: "部门",
	4: "小组",
}

// 菜单类型
var MENU_TYPE_LIST = map[int]string{
	0: "菜单",
	1: "节点",
}

// 城市等级
var CITY_LEVEL = map[int]string{
	1: "省份",
	2: "城市",
	3: "县区",
	4: "街道",
}

// 配置项类型
var CONFIG_DATA_TYPE_LIST = map[string]string{
	"text":     "单行文本",
	"textarea": "多行文本",
	"ueditor":  "富文本编辑器",
	"date":     "日期",
	"datetime": "时间",
	"number":   "数字",
	"select":   "下拉框",
	"radio":    "单选框",
	"checkbox": "复选框",
	"image":    "单张图片",
	"images":   "多张图片",
	"password": "密码",
	"icon":     "字体图标",
	"file":     "单个文件",
	"files":    "多个文件",
	"hidden":   "隐藏",
	"readonly": "只读文本",
}

// 友链类型
var LINK_TYPE_LIST = map[int]string{
	1: "友情链接",
	2: "合作伙伴",
}

// 友链形式
var LINK_FORM_LIST = map[int]string{
	1: "文字链接",
	2: "图片链接",
}

// 友链平台
var LINK_PLATFORM_LIST = map[int]string{
	1: "PC站",
	2: "WAP站",
	3: "小程序",
	4: "APP应用",
}

// 站点类型
var ITEM_TYPE_LIST = map[int]string{
	1: "国内站点",
	2: "国外站点",
	3: "其他站点",
}

// 广告位所属平台
var ADSORT_PLATFORM_LIST = map[int]string{
	1: "PC站",
	2: "WAP站",
	3: "小程序",
	4: "APP应用",
}

// 广告类型
var AD_TYPE_LIST = map[int]string{
	1: "图片",
	2: "文字",
	3: "视频",
	4: "其他",
}

// 通知来源
var NOTICE_SOURCE_LIST = map[int]string{
	1: "内部通知",
	2: "外部通知",
}

// 会员设备类型
var MEMBER_DEVICE_LIST = map[int]string{
	1: "苹果",
	2: "安卓",
	3: "WAP站",
	4: "PC站",
	5: "后台添加",
}

// 会员来源
var MEMBER_SOURCE_LIST = map[int]string{
	1: "注册会员",
	2: "马甲会员",
}
