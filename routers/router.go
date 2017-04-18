package routers

import (
	"gomage/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/admin/sys", &controllers.AdminController{}, "get,post:System")
	beego.Router("/admin/getsys", &controllers.AdminController{}, "get:GetSysById")
	beego.Router("/admin/style", &controllers.AdminController{}, "get,post:Style")
	beego.Router("/admin/", &controllers.AdminController{}, "get:Index")
	beego.Router("/admin/clean", &controllers.AdminController{}, "get:CleanUp")
	beego.Router("/admin/index", &controllers.AdminController{}, "get:Index")
	beego.Router("/admin/del", &controllers.AdminController{}, "get:Del")
	beego.Router("/admin/update", &controllers.AdminController{}, "get:Update")
	beego.Router("/admin/update_pwd", &controllers.AdminController{}, "post:UpdatePwd")
	beego.Router("/admin/preview", &controllers.AdminController{}, "get:Preview") //图片处理预览
	beego.Router("/admin/import", &controllers.AdminController{}, "post:Import")  //导入样式
	beego.Router("/admin/login", &controllers.LoginController{}, "get,post:Login")
	beego.Router("/admin/logout", &controllers.AdminController{}, "get,post:Logout")
	beego.Router("/*", &controllers.MainController{}, "get:View")

}
