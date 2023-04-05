package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-frame/entity/cmsdto"
	"reflect"
)

type sDbNewCms struct {
	db     gdb.DB
	prefix string
}

func DbNewCms() *sDbNewCms {
	return &sDbNewCms{db: g.DB("newcms"), prefix: "cms_"}
}

// 询条件
//type AppApiReq struct {
//	TableName string `p:"TableName"` // 表明
//	Page     int    `p:"page"`     // 页码
//	Limit    int    `p:"limit"`    // 每页数
//	Where    string  `p:"where"`	//查询条件 json字符串 内部查询转换g.map结构体
//}
func (cms *sDbNewCms) ReadByParam(option *cmsdto.ListReq) (list interface{}, count int, err error) {
	var tableName = cms.prefix + option.TableName
	var pageSize = option.Limit
	var pageNow = option.Page
	var limitOne = (pageNow - 1) * pageSize

	var where g.Map = g.Map{}
	_ = json.Unmarshal([]byte(option.Where), &where)
	var res gdb.Result

	count, err = cms.db.Model(tableName).Where(where).Count()
	var order = fmt.Sprintf("%v", option.Order)

	res, err = cms.db.Model(tableName).Where(where).Order(order).Limit(limitOne, pageSize).All()

	if err != nil {
		fmt.Printf("err->%v\r\n", err)
	}
	list = res
	return
}

//返回值为错误信息和更新成功行数
func (cms *sDbNewCms) UpByParam(option *cmsdto.UpdateReq, UserId int) (rows int64, err error) {
	var tableName = cms.prefix + option.TableName
	var where g.Map = g.Map{}
	_ = json.Unmarshal([]byte(option.Where), &where)

	var data g.Map = g.Map{}
	err = json.Unmarshal([]byte(option.Data), &data)
	if err != nil {
		return
	}
	data["update_time"] = gtime.Now()
	data["update_user"] = UserId
	var res sql.Result
	res, err = cms.db.Model(tableName).Where(where).Update(data)
	if err != nil {
		return 0, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return
}

func (cms *sDbNewCms) AddByParam(option *cmsdto.AddReq, UserId int) (rows int64, add_id int64, err error) {
	var tableName = cms.prefix + option.TableName
	var data g.List = g.List{}
	err = json.Unmarshal([]byte(option.Data), &data)
	if err != nil {
		return
	}
	fmt.Printf("addJson:%v", data)
	var res sql.Result
	for i := 0; i < len(data); i++ {
		data[i]["create_time"] = gtime.Now()
		data[i]["create_user"] = UserId
	}
	res, err = cms.db.Model(tableName).Save(data)
	rows, err = res.RowsAffected()
	add_id, err = res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return
}

func (cms *sDbNewCms) DelByParam(option *cmsdto.DelReq) (rows int64, err error) {
	var tableName = cms.prefix + option.TableName
	var where g.Map = g.Map{}
	err = json.Unmarshal([]byte(option.Where), &where)
	if err != nil {
		return
	}
	//fmt.Printf("delewhere->%v",where)
	var res sql.Result
	res, err = cms.db.Model(tableName).Where(where).Delete()
	rows, err = res.RowsAffected()

	if err != nil {
		return 0, err
	}

	go func() {
		//开新线程执行不影响正常返回
		for {
			//反射调用方法
			in := make([]reflect.Value, 1)
			var delId interface{} = 0
			switch option.TableName {
			case "capp":
				delId = where["appid"]
				break
			}
			DoMethodName := "Del_" + option.TableName
			t := reflect.TypeOf(cms)
			in[0] = reflect.ValueOf(delId)
			for i := 0; i < t.NumMethod(); i++ {
				method := t.Method(i)
				if method.Name == DoMethodName {
					_ = reflect.ValueOf(cms).MethodByName(DoMethodName).Call(in)
					break
				}
			}
			break
		}

	}()

	return
}

func (cms *sDbNewCms) UpVideoStatus(where g.Map) (err error) {
	var tableName = cms.prefix + "video"
	_, err = cms.db.Model(tableName).Where(where).Save(g.Map{"status": 1})
	return
}
