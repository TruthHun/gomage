package controllers

import (
	"gomage/helper"
	"gomage/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/disintegration/imaging"
	//"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"fmt"

	"image"

	"github.com/astaxie/beego/cache"
)

type MainController struct {
	BaseController
	//beego.Controller
}

//图片浏览
func (this *MainController) View() {
	var (
		date, md5str, cache_file, host, xhost, referer string
		styles                                         []models.Style
		style                                          models.Style
		img                                            image.Image
	)
	referer = this.Ctx.Request.Header.Get("Referer")
	host = this.Ctx.Request.Host
	xhost = this.Ctx.Request.Header.Get("X-Forwarded-Host")
	if len(xhost) > 0 {
		host = xhost
	}
	slice := strings.Split(host, ":")
	if len(slice) > 1 {
		//带有端口的去除端口
		host = slice[0]
	}

	//根据域名查询系统中站点的配置
	sys := models.GetSysByHost(host)
	//查询不到数据或者域名状态为false，返回域名不存在
	if sys.Id == 0 || sys.Status == false {
		this.Ctx.WriteString(I18n(this.Lang, "domain_disabled"))
		return
	}

	//检测referer
	sys.Referer = strings.TrimSpace(sys.Referer)
	if len(sys.Referer) > 0 {
		fmt.Println(sys.Referer, referer)
		referer_slice := strings.Split(sys.Referer, ",")
		allowed := false
		for _, v := range referer_slice {
			if strings.Contains(referer, v) || strings.Contains(host, v) {
				allowed = true
				break
			}
		}
		if allowed == false {
			this.Ctx.WriteString("Forbidden referer")
			return
		}
	}
	splat := this.GetString(":splat")
	//是否是获取图片信息
	IsGetInfo := strings.Contains(splat, "@@info")
	styles, err := CacheStyles(host, sys.Id)
	style_len := len(styles)
	if (style_len == 0 || err != nil) && !IsGetInfo {
		if err != nil {
			this.Ctx.WriteString(err.Error())
		} else {
			this.Ctx.WriteString(host + ":" + I18n(this.Lang, "style_disabled"))
		}
		return
	}

	//favicon图片直接输出显示
	if strings.HasSuffix(splat, ".ico") {
		b, _ := ioutil.ReadFile("./" + strings.TrimLeft(splat, "/"))
		this.Ctx.ResponseWriter.Write(b)
	} else {

		//如果是获取图片信息，则返回图片信息过去
		if IsGetInfo {
			dst := strings.TrimRight(sys.Prefix, "/") + "/" + strings.TrimLeft(strings.TrimSuffix(splat, "@@info"), "/")
			info, err := helper.ImageInfo(dst)
			if err != nil {
				this.Ctx.WriteString(err.Error())
			} else {
				this.Data["json"] = info
				this.ServeJSON()
			}
			return
		}
		writer := this.Ctx.ResponseWriter
		//切割图片样式，返回图片路径以及样式别名
		SplitStyle := func(path string) (string, string) {
			slice := strings.Split(path, sys.Segment)
			rule := strings.TrimSpace(slice[len(slice)-1])
			path = strings.TrimSuffix(path, sys.Segment+rule)
			return path, rule
		}
		path, rule := SplitStyle(splat)
		//图片文件路径组装
		sys.Prefix = strings.TrimSpace(sys.Prefix)
		if len(sys.Prefix) == 0 {
			sys.Prefix = "." //表示当前同级目录
		}
		path = strings.TrimRight(sys.Prefix, "/") + "/" + strings.TrimLeft(path, "/")

		for _, v := range styles {
			if v.Rule == rule {
				if v.Status {
					style = v
					break //终止循环
				}
			}
		}

		//图片处理样式存在，则进行处理
		if style.Id > 0 {
			//流程：
			//1、看下上次的请求时间，如果浏览器有缓存，则返回304
			//2、浏览器没有缓存，查看服务器是否存在图片缓存，有则直接读取输出
			//3、裁剪图片，如果设置了缓存，则缓存图片，否则不缓存

			//获取图片文件信息

			_, err := os.Stat(path)
			if err != nil {
				this.Ctx.WriteString("No such file or directory")
				return
			}

			//打开图片文件
			img, err = imaging.Open(path)
			if err != nil {
				this.Ctx.WriteString(err.Error())
				return
			}
			date = time.Unix(int64(style.Time), 0).Format(time.FixedZone(time.RFC1123, 0).String())
			date = strings.Replace(date, "CST", "GMT", -1) //替换成GMT
			//时间校验
			since, ok := this.Ctx.Request.Header["If-Modified-Since"]
			if ok && since[0] == date {
				this.Ctx.ResponseWriter.WriteHeader(http.StatusNotModified)
				return
			}
			//上次修改时间
			writer.Header().Add("Last-Modified", date)
			//根据参数查找下是否存在缓存文件
			style_str := fmt.Sprintf("%v-%v", style.Rule, path)
			md5str = helper.MyMD5(style_str)

			//获取图片格式数字以及图片格式后缀
			GetFormat := func(ext string) (imaging.Format, string) {
				var f imaging.Format
				var e string
				//图片输出格式
				switch ext {
				case "jpeg", "jpg":
					f, e = imaging.JPEG, "jpg"
				case "png":
					f, e = imaging.PNG, "png"
				case "gif":
					f, e = imaging.GIF, "gif"
				case "bmp":
					f, e = imaging.BMP, "bmp"
				case "tiff", "tif":
					f, e = imaging.TIFF, "tif"
				default:
					f, e = imaging.Format(-1), ""
				}
				return f, e
			}
			format, cache_ext := GetFormat(style.Ext)
			if cache_ext == "" {
				//原图格式输出
				slice := strings.Split(path, ".")
				format, cache_ext = GetFormat(strings.ToLower(strings.TrimSpace(slice[len(slice)-1])))
			}
			md5slice := strings.Split(md5str, "")
			cache_folder := "./cache/" + host + "/" + style.Rule + "/" + strings.Join(md5slice[0:5], "/") + "/"
			cache_file = cache_folder + md5str + "." + cache_ext

			//如果开启了缓存
			if sys.IsCache {
				info, err := os.Stat(cache_file)
				//如果缓存文件存在且缓存文件的生成时间大于样式的生成时间，则直接读取缓存文件返回
				if err == nil && int(info.ModTime().Unix()) > style.Time {
					//以原图的Last-Modified为准
					b, _ := ioutil.ReadFile(cache_file)
					writer.Write(b)
					return
				}
			}
			//调用图片处理
			img = helper.Proccess(img, style.Width, style.Height, style.Method, style.IsZoom)
			//图片水印处理
			style.Watermark = strings.TrimSpace(style.Watermark)
			if len(style.Watermark) > 0 {
				imgwatermark, err := helper.AddWatermark(img, style.Watermark, style.WatermarkPosition, style.Top, style.Right, style.Bottom, style.Left)
				if err == nil {
					img = imgwatermark
				}
			}
			//生成图片缓存
			if sys.IsCache {
				os.MkdirAll(cache_folder, 0777)
				imaging.Save(img, cache_file)
			}
			imaging.Encode(writer, img, format)
		} else {
			//如果规则不存在，且没有进行原图保护，则直接显示原图
			if sys.Protected == false {
				//路径再处理，是因为当规则分隔符是“/”的时候，图片文件路径就会出现不正确的情况，
				//如图片路为./uploads/1.jpg，而分隔符是“/”，则path就变成了./uploads，并不是图片

				path = strings.TrimRight(path, sys.Segment)
				helper.HCPrint(path)
				b, err := ioutil.ReadFile(path)
				helper.HCPrint(path, err)
				if err == nil {
					info, _ := os.Stat(path)
					if helper.IsImage(info.Name()) {
						//文件更改时间
						date = strings.Replace(info.ModTime().Format(time.FixedZone(time.RFC1123, 0).String()), "CST", "GMT", -1)
						//时间校验
						since, ok := this.Ctx.Request.Header["If-Modified-Since"]
						if ok && since[0] == date {
							this.Ctx.ResponseWriter.WriteHeader(http.StatusNotModified)
						} else {
							writer.Header().Add("Last-Modified", date)
							writer.Write(b)
						}
					} else {
						this.Ctx.WriteString("The image format does not meet the requirements")
					}
				} else {
					this.Ctx.WriteString("No such file or directory")
				}
			} else {
				this.Ctx.WriteString("Do not access the original image")
			}
		}
	}
}

//语言函数
func I18n(lang, tag string) string {
	key := fmt.Sprintf("%v::%v", lang, tag)
	v := beego.AppConfig.String(key)
	if len(v) == 0 {
		v = fmt.Sprintf("{{%s}}", tag)
	}
	return v
}

//返回选项
func MethodName(option int, lang string) string {
	data := map[int]string{
		11: I18n(lang, "method11"),
		12: I18n(lang, "method12"),
		21: I18n(lang, "method21"),
		22: I18n(lang, "method22"),
	}
	return data[option]
}

//根据host查询缓存的样式，如果缓存的样式不存在，则查询数据库数据，否则返回缓存数据。
func CacheStyles(host string, sid int) ([]models.Style, error) {
	var styles []models.Style
	var ok bool
	var err error
	bc, err := cache.NewCache("file", `{"CachePath":"./cache/runtime","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":5}`) //beego cache，beego缓存；过期时间5秒
	if err != nil {
		return styles, err
	}
	//判断当前域名的缓存是否存在
	if bc.IsExist(host) {
		styles, ok = bc.Get(host).([]models.Style)
		if ok && len(styles) > 0 {
			return styles, nil
		}
	}
	//先根据host查询样式缓存，如果样式缓存不存在，则再从数据库中查询样式，并缓存新查询的样式
	styles = models.GetStyleDataBySid(sid)
	if len(styles) > 0 {
		err = bc.Put(host, styles, 5*time.Second)
	}
	return styles, err
}
