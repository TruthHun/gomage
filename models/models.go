package models

import (
	"fmt"
	"wangpan/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//初始化数据库注册
func init() {
	//初始化数据库
	RegisterDB()
	runmode := beego.AppConfig.String("runmode")
	if runmode == "prod" {
		orm.Debug = false
		orm.RunSyncdb("default", false, false)
	} else {
		orm.Debug = true
		orm.RunSyncdb("default", false, true)
	}
}

//注册数据库
func RegisterDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	models := []interface{}{
		new(Style),
		new(Admin),
		new(System),
	}
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db_prefix"), models...)
	db_user := beego.AppConfig.String("db_user")
	db_password := beego.AppConfig.String("db_password")
	db_database := beego.AppConfig.String("db_database")
	db_charset := beego.AppConfig.String("db_charset")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	dblink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", db_user, db_password, db_host, db_port, db_database, db_charset)

	//下面两个参数后面要放到app.conf提供用户配置使用
	// (可选)设置最大空闲连接
	maxIdle := 50
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := 50
	orm.RegisterDataBase("default", "mysql", dblink, maxIdle, maxConn)
}

//管理员登录
func Login(username, password string) int {
	var admin Admin
	table := GetTable("admin")
	o := orm.NewOrm()
	password = tools.MyMD5(password)
	o.QueryTable(table).Filter("username", username).Filter("password", password).One(&admin)
	return admin.Id
}

//获取数据表
func GetTable(tab string) string {
	return beego.AppConfig.DefaultString("db_prefix", "hc_") + tab
}

//根据表更新字段内容
func Update(table string, field string, value interface{}, id ...int) (int, error) {
	i, err := orm.NewOrm().QueryTable(GetTable(table)).Filter("Id__in", id).Update(orm.Params{
		field: value,
	})
	return int(i), err
}

//新增或更新数据
func AddSite(sys System) (int64, error) {
	o := orm.NewOrm()
	if sys.Id > 0 {
		fields := []string{"Sitename", "Host", "Prefix", "Protected", "Segment", "IsCache"}
		return o.Update(&sys, fields...)
	} else {
		sys.Status = true
		return o.Insert(&sys)
	}
}

//系统设置的数据
func GetSysData() []orm.Params {
	var sys []orm.Params
	orm.NewOrm().QueryTable(GetTable("system")).OrderBy("-Id").Values(&sys)
	return sys
}

//获取站点id获取站点的样式数据
func GetStyleDataBySid(sid int) []Style {
	var style []Style
	orm.NewOrm().QueryTable(GetTable("style")).Filter("Sid", sid).OrderBy("-id").All(&style)
	return style
}

//获取图片处理样式
func GetStyleByRule(rule string) Style {
	var style Style
	orm.NewOrm().QueryTable(GetTable("style")).Filter("rule", rule).One(&style)
	return style
}

//根据id获取样式数据
func GetStyleById(id int) Style {
	var style Style
	orm.NewOrm().QueryTable(GetTable("style")).Filter("id", id).One(&style)
	return style
}

//根据host获取站点样式
func GetStyleByHost(host string) []Style {
	var (
		sys System
	)
	o := orm.NewOrm()
	o.QueryTable(GetTable("system")).Filter("host", host).Filter("status", 1).One(&sys)
	return GetStyleDataBySid(sys.Id)
}

//根据id获取站点信息
//id为0时获取最新的一条记录
func GetSysById(id int) System {
	var sys System
	qs := orm.NewOrm().QueryTable(GetTable("system"))
	if id == 0 {
		qs.OrderBy("-Id").One(&sys)
	} else {
		qs.Filter("Id", id).One(&sys)
	}
	return sys
}

//根据host查询站点配置
func GetSysByHost(host string) System {
	var sys System
	orm.NewOrm().QueryTable(GetTable("system")).Filter("host", host).One(&sys)
	return sys
}

//根据id删除样式
func DelStyleById(id int) (bool, error) {
	i, err := orm.NewOrm().QueryTable(GetTable("style")).Filter("id", id).Delete()
	if i > 0 {
		return true, err
	}
	return false, err
}

//根据id删除指定表的内容
func DelById(table string, id ...interface{}) (int64, error) {
	return orm.NewOrm().QueryTable(GetTable(table)).Filter("Id__in", id...).Delete()
}

//根据id更新指定表的字段值
func UpdateById(table, field string, value interface{}, id ...interface{}) (int64, error) {
	return orm.NewOrm().QueryTable(GetTable(table)).Filter("Id__in", id...).Update(orm.Params{field: value})
}

//根据条件删除指定表的内容
func DelByCond(table string, cond *orm.Condition) (int64, error) {
	return orm.NewOrm().QueryTable(GetTable(table)).SetCond(cond).Delete()
}

//根据id获取管理员信息
func GetAdmin(id int) Admin {
	var admin Admin
	orm.NewOrm().QueryTable(GetTable("admin")).Filter("id", id).One(&admin)
	return admin
}
