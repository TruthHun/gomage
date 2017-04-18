package controllers

import (
	"fmt"
	"gomage/models"
)

type LoginController struct {
	BaseController
}

//登录页面
func (this *LoginController) Login() {
	fmt.Println("登录页面")
	login, ok := this.GetSession("uid").(int)
	//如果已经登录，则天转到后台
	if ok && login > 0 {
		this.Redirect("/admin/index", 302)
	}
	if this.Ctx.Request.Method == "POST" {
		uid := models.Login(this.GetString("username"), this.GetString("password"))
		if uid > 0 {
			//登录成功，设置session
			this.SetSession("uid", uid)
			this.Redirect("/admin/index", 302)
		} else {
			//账号或密码不正确
			this.Redirect("/admin/login?msg="+I18n(this.Lang, "error_login"), 302)
		}
		//post方式
	} else {
		//get方式
		this.Data["Msg"] = this.GetString("msg")
		this.TplName = "login.html"
	}
}
