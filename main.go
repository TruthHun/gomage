package main

import (
	"gomage/controllers"
	_ "gomage/routers"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
)

//初始化模板函数
func init() {

	beego.AddFuncMap("Option", controllers.MethodName)

	beego.AddFuncMap("i18n", controllers.I18n)

}
func main() {
	beego.Run()
}
