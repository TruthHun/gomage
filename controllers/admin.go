package controllers

import (
	"gomage/helper"
	"gomage/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//"strings"
	"fmt"
	"regexp"

	"encoding/json"

	"io"

	"os"

	"github.com/disintegration/imaging"
)

type BaseController struct {
	beego.Controller
	IsLogin int
	Lang    string
	Host    string
}

//初始化
func (this *BaseController) Prepare() {
	lang := this.Ctx.GetCookie("lang")
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			lang = al[:5]
		}
	}
	if lang != "zh-CN" {
		this.Lang = "en-US"
	} else {
		this.Lang = lang
	}
	this.Data["Lang"] = this.Lang
	runmode := beego.AppConfig.String("runmode")
	if runmode != "prod" {
		this.Data["StaticVersion"] = time.Now().Unix()
	} else {
		this.Data["StaticVersion"] = "gomage"
	}
}

type AdminController struct {
	BaseController
}

//构造函数
func (this *AdminController) Prepare() {
	this.BaseController.Prepare()
	uid, ok := this.GetSession("uid").(int)
	if !ok || uid == 0 {
		this.Redirect("/admin/login", 302)
	}
	this.IsLogin = uid
	this.Layout = "layout.html"
	this.Data["TplStatic"] = "/static"
	this.Data["Uid"] = uid
}

//后台首页
func (this *AdminController) Index() {

	//语言设置
	lang := this.GetString("lang", "")
	if len(lang) > 0 {
		if lang == "cn" {
			this.Ctx.SetCookie("lang", "zh-CN")
		} else {
			this.Ctx.SetCookie("lang", "en-US")
		}
		data := map[string]interface{}{"status": 1, "msg": I18n(this.Lang, "update_success")}
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	this.Data["Title"] = I18n(this.Lang, "dashboard")
	this.Data["IsIndex"] = true
	this.TplName = "index.html"
}

//后台系统设置
func (this *AdminController) System() {
	if this.Ctx.Request.Method == "POST" {
		//新增或者更新站点
		sys := models.System{}
		this.ParseForm(&sys)
		i, err := models.AddSite(sys)
		var data = map[string]interface{}{"status": 0}
		if sys.Id > 0 {
			if i > 0 && err == nil {
				data["status"] = 1
				data["msg"] = I18n(this.Lang, "update_success")
			} else {
				data["msg"] = I18n(this.Lang, "update_fail")
			}
		} else {
			if i > 0 && err == nil {
				data["status"] = 1
				data["msg"] = I18n(this.Lang, "create_success")
			} else {
				data["msg"] = I18n(this.Lang, "fail_create_site")
			}
		}
		this.Data["json"] = data
		this.ServeJSON()
	} else {
		this.Data["IsSys"] = true
		this.Data["Title"] = I18n(this.Lang, "site_config")
		action := this.GetString("action")
		if action == "edit" {
			id, _ := this.GetInt("id")
			if id > 0 {
				this.Data["Sys"] = models.GetSysById(id)
				this.TplName = "site_setting.html"
			} else {
				this.Redirect("/admin/sys", 302)
			}
		} else {
			this.Data["Sys"] = models.GetSysData()
			this.TplName = "site.html"
		}

	}
}

//后台样式管理
func (this *AdminController) Style() {
	this.Data["IsStyle"] = true
	//创建或更新样式
	if this.Ctx.Request.Method == "POST" {
		var (
			style     models.Style
			flag      bool = true
			i         int64
			err       error
			msg       string
			is_update bool = false
		)
		this.ParseForm(&style)
		data := map[string]interface{}{"status": 0, "msg": I18n(this.Lang, "create_fail")}
		if style.Id > 0 {
			is_update = true
			msg = I18n(this.Lang, "update_fail")
		}
		//宽度和高度均不能小于0
		if style.Width < 0 || style.Height < 0 {
			msg = I18n(this.Lang, "img_w_h_lt_0_err")
			flag = false
		}

		//如果裁剪方式是21、22，则宽或高可以为0，否则不能为0
		if flag && style.Method < 20 && (style.Width == 0 || style.Height == 0) {
			flag = false
			msg = I18n(this.Lang, "img_w_h_err")
		}
		if flag {
			if style.Method == 21 && style.Width == 0 {
				flag = false
				msg = I18n(this.Lang, "img_w_err")
			}
			if style.Method == 22 && style.Height == 0 {
				flag = false
				msg = I18n(this.Lang, "img_h_err")
			}
		}

		//验证样式规则名称，样式规则名称仅限字母、数字、下划线(_)、短横线(-)以及小数点
		if flag {
			patern := fmt.Sprintf(`[a-zA-z0-9\.\-\_]{%v}`, len(style.Rule))
			b, err := regexp.MatchString(patern, style.Rule)
			if b == false || err != nil {
				msg = I18n(this.Lang, "rule_tips")
				flag = false
			}
		}

		if flag {
			style.Time = int(time.Now().Unix())
			if style.Id > 0 {
				i, err = orm.NewOrm().Update(&style)
			} else {
				i, err = orm.NewOrm().Insert(&style)
			}
			//更新或者插入之后，style的id都大于0
			if i > 0 {
				msg = I18n(this.Lang, "create_success")
				if is_update == true {
					msg = I18n(this.Lang, "update_success")
					oldRule := this.GetString("oldRule")

					if oldRule != style.Rule {
						sys := models.GetSysById(style.Sid)
						go os.RemoveAll("./cache/" + sys.Host + "/" + oldRule)
					}
				}
				data["status"] = 1
			} else {
				if err == nil {
					if is_update {
						data["msg"] = I18n(this.Lang, "update_fail")
					} else {
						data["msg"] = I18n(this.Lang, "create_fail")
					}
				} else {
					data["msg"] = err.Error()
				}

			}
		}
		data["msg"] = msg
		this.Data["json"] = data
		this.ServeJSON()
	} else {
		action := this.GetString("action")
		this.Data["Action"] = action
		sid, _ := this.GetInt("sid")
		options := models.GetSysData()
		if sid == 0 {
			if len(options) > 0 {
				sid = int(options[0]["Id"].(int64))
			}
		}
		this.Data["Sid"] = sid
		this.Data["Options"] = options

		switch action {
		case "add":
			this.Data["Title"] = I18n(this.Lang, "style_create")
			this.Data["Action"] = "add"
			this.TplName = "style_add.html"
		case "edit":
			id, _ := this.GetInt("id")
			if id == 0 {
				this.Redirect("/admin/style?error=id error", 302)
			} else {
				style := models.GetStyleById(id)
				if style.Id == 0 {
					this.Redirect("/admin/style?error=style error", 302)
					return
				}
				this.Data["Style"] = style
				this.Data["Title"] = I18n(this.Lang, "style_edit")
				this.Data["Action"] = "edit"
				this.TplName = "style_edit.html"
			}
		case "export":
			data := models.GetStyleDataBySid(sid)
			b, _ := json.Marshal(data)
			this.Ctx.ResponseWriter.Header().Set("Content-disposition", "attachment; filename=style.json")
			this.Ctx.ResponseWriter.Write(b)
			return
		default:
			this.Data["Title"] = I18n(this.Lang, "style_list")
			this.Data["Data"] = models.GetStyleDataBySid(sid)
			this.Data["Action"] = "list"
			this.TplName = "style.html"
		}

	}
}

//导入样式
func (this *AdminController) Import() {
	var (
		styles []models.Style
		err    error
		n, sid int
	)
	data := map[string]interface{}{"status": 0, "msg": I18n(this.Lang, "style_import_err")}
	sid, _ = this.GetInt("sid")
	if sid > 0 {
		f, _, _ := this.GetFile("style")
		defer f.Close()
		chunks := make([]byte, 1, 1024)
		buf := make([]byte, 1)
		for {
			n, err = f.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if 0 == n {
				break
			}
			chunks = append(chunks, buf[:n]...)
			//fmt.Println(string(buf[:n]))
		}
		err = json.Unmarshal(chunks[1:], &styles)
		if err == nil {
			if len(styles) > 0 {
				o := orm.NewOrm()
				for _, style := range styles {
					style.Id = 0
					style.Sid = sid
					style.Status = true
					style.Time = int(time.Now().Unix())
					o.Insert(&style)
				}
				data["status"] = 1
				data["msg"] = I18n(this.Lang, "style_import_success")
			} else {
				data["msg"] = I18n(this.Lang, "style_data_err")
			}
		} else {
			data["msg"] = err.Error()
		}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

//根据id更新字段内容
func (this *AdminController) Update() {
	data := map[string]interface{}{
		"status": 0,
		"msg":    I18n(this.Lang, "parameter_incorrect"),
	}
	id, _ := this.GetInt("id")
	table := this.GetString("table")
	field := this.GetString("field")
	value := this.GetString("value")
	if id > 0 {
		i, err := models.Update(table, field, value, id)
		if i > 0 {
			data["msg"] = I18n(this.Lang, "update_success")
			data["status"] = 1
		} else {
			if err == nil {
				data["msg"] = I18n(this.Lang, "update_failed")
			} else {
				data["msg"] = err.Error()
			}
		}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

//根据id获取数据
func (this *AdminController) GetSysById() {
	id, _ := this.GetInt("id")
	this.Data["json"] = models.GetSysById(id)
	this.ServeJSON()
}

//删除内容
func (this *AdminController) Del() {
	var sys models.System
	data := map[string]interface{}{
		"msg":    I18n(this.Lang, "parameter_incorrect"),
		"status": 0,
	}
	id, _ := this.GetInt("id")
	table := this.GetString("table")
	if id > 0 {
		if table == "system" {
			sys = models.GetSysById(id)
		}
		i, err := models.DelById(table, id)
		if i > 0 {
			if table == "system" {
				cond := orm.NewCondition()
				models.DelByCond("style", cond.And("Sid", id))
				//删除成功之后，还要去删除缓存文件等
				go os.RemoveAll("./cache/" + sys.Host)
			}
			data["msg"] = I18n(this.Lang, "del_success")
			data["status"] = 1
		} else {
			if err == nil {
				data["msg"] = I18n(this.Lang, "del_not_exist_err")
			} else {
				data["msg"] = err.Error()
			}
			data["status"] = 0
		}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

//图片处理预览
func (this *AdminController) Preview() {
	img, err := imaging.Open("./static/images/example.jpg") //预览时需要处理的图片
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
	waterposition, _ := this.GetInt("waterposition", 9) //水印位置
	waterpath := this.GetString("waterpath", "")        //水印路径
	width, _ := this.GetInt("width")                    //图片宽度
	height, _ := this.GetInt("height")                  //图片高度
	method, _ := this.GetInt("method")                  //图片处理方式
	zoom, _ := this.GetBool("zoom", true)               //是否允许图片放大
	ext := this.GetString("ext", "-")                   //图片格式
	top, _ := this.GetInt("top")
	right, _ := this.GetInt("right")
	left, _ := this.GetInt("left")
	bottom, _ := this.GetInt("bottom")
	//显示原图
	if (width == 0 && height == 0) || method == 0 {
		imaging.Encode(this.Ctx.ResponseWriter, img, imaging.JPEG)
		return
	}
	img = helper.Proccess(img, width, height, method, zoom)
	if len(waterpath) > 0 {
		waterimg, err := helper.AddWatermark(img, waterpath, waterposition, top, right, bottom, left)
		if err == nil {
			img = waterimg
		}
	}
	var f imaging.Format
	switch ext {
	case "jpeg", "jpg":
		f = imaging.JPEG
	case "png":
		f = imaging.PNG
	case "gif":
		f = imaging.GIF
	case "bmp":
		f = imaging.BMP
	case "tiff", "tif":
		f = imaging.TIFF
	default:
		f = imaging.JPEG
	}
	imaging.Encode(this.Ctx.ResponseWriter, img, f)
}

//清理缓存文件
func (this *AdminController) CleanUp() {
	var data = map[string]interface{}{"status": 1, "msg": I18n(this.Lang, "tips_clear_cache")}
	sid, _ := this.GetInt("sid")
	if sid > 0 {
		sys := models.GetSysById(sid)
		go os.RemoveAll("./cache/" + sys.Host)
	}
	this.Data["json"] = data
	this.ServeJSON()
}

//更新账号和密码
func (this *AdminController) UpdatePwd() {
	data := map[string]interface{}{"status": 0, "msg": I18n(this.Lang, "parameter_incorrect")}
	pwd_old := this.GetString("password_old")
	pwd_new := this.GetString("password_new")
	pwd_ensure := this.GetString("password_ensure")
	uid, _ := this.GetInt("uid")
	if uid > 0 {
		if pwd_new != pwd_old && pwd_ensure == pwd_new {
			admin := models.GetAdmin(uid)
			if admin.Password == helper.MyMD5(pwd_old) {
				pwd := helper.MyMD5(pwd_new)
				i, err := models.Update("admin", "password", pwd, uid)
				if i > 0 {
					data["status"] = 1
					data["msg"] = I18n(this.Lang, "update_success")
				} else {
					if err == nil {
						data["msg"] = err.Error()
					} else {
						data["msg"] = I18n(this.Lang, "fail_pwd")
					}
				}
			} else {
				data["msg"] = I18n(this.Lang, "ori_pwd_err")
			}
		} else {
			data["msg"] = I18n(this.Lang, "ori_pwd_eq_new_pwd_err")
		}
	}
	this.Data["json"] = data
	this.ServeJSON()
}

//退出登录
func (this *AdminController) Logout() {
	this.DestroySession()
	this.Redirect("/admin/login", 302)
}
