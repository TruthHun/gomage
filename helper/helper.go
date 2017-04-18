package helper

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"net/mail"

	"github.com/alexcesaro/mail/mailer"
)

//时间戳格式化
func TimestampFormat(timestamp int, format string) string {
	return time.Unix(int64(timestamp), 0).Format(format)
}

//转化成小写
func StrToLower(str string) string {
	return strings.ToLower(str)
}

//转化成整型
func ParseInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//分页函数
//rollPage:展示分页的个数
//totalRows：总记录
//currentPage:每页显示记录数
//urlPrefix:url链接前缀
//urlParams:url键值对参数
func Paginations(rollPage, totalRows, listRows, currentPage int, urlPrefix string, urlParams ...interface{}) string {
	var (
		htmlPage, path string
		pages          []int
		params         []string
	)
	//总页数
	totalPage := totalRows / listRows
	if totalRows%listRows > 0 {
		totalPage += 1
	}
	//只有1页的时候，不分页
	if totalPage < 2 {
		return ""
	}
	params_len := len(urlParams)
	if params_len > 0 {
		if params_len%2 > 0 {
			params_len = params_len - 1
		}
		for i := 0; i < params_len; {
			key := strings.TrimSpace(fmt.Sprintf("%v", urlParams[i]))
			val := strings.TrimSpace(fmt.Sprintf("%v", urlParams[i+1]))
			//键存在，同时值不为0也不为空
			if len(key) > 0 && len(val) > 0 && val != "0" {
				params = append(params, key, val)
			}
			i = i + 2
		}
	}

	path = strings.Trim(urlPrefix, "/")
	if len(params) > 0 {
		path = path + "/" + strings.Trim(strings.Join(params, "/"), "/")
	}
	//最后再处理一次“/”，是为了防止urlPrifix参数为空时，出现多余的“/”
	path = "/" + strings.Trim(path, "/") + "/p/"

	if currentPage > totalPage {
		currentPage = totalPage
	}
	if currentPage < 1 {
		currentPage = 1
	}
	index := 0
	rp := rollPage * 2
	for i := rp; i > 0; i-- {
		p := currentPage + rollPage - i
		if p > 0 && p <= totalPage {

			pages = append(pages, p)
		}
	}
	for k, v := range pages {
		if v == currentPage {
			index = k
		}
	}
	pages_len := len(pages)
	if currentPage > 1 {
		htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`1">1..</a></li><li><a class="num" href="`+path+`%d"><<</a></li>`, currentPage-1)
	}
	if pages_len <= rollPage {
		for _, v := range pages {
			if v == currentPage {
				htmlPage += fmt.Sprintf(`<li class="active"><a href="javascript:void(0);">%d</a></li>`, v)
			} else {
				htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`%d">%d</a></li>`, v, v)
			}
		}

	} else {
		index_min := index - rollPage/2
		index_max := index + rollPage/2
		page_slice := make([]int, 0)
		if index_min > 0 && index_max < pages_len { //切片索引未越界
			page_slice = pages[index_min:index_max]
		} else {
			if index_min < 0 {
				page_slice = pages[0:rollPage]
			} else if index_max > pages_len {
				page_slice = pages[(pages_len - rollPage):pages_len]
			} else {
				page_slice = pages[index_min:index_max]
			}

		}

		for _, v := range page_slice {
			if v == currentPage {
				htmlPage += fmt.Sprintf(`<li class="active"><a href="javascript:void(0);">%d</a></li>`, v)
			} else {
				htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`%d">%d</a></li>`, v, v)
			}
		}

	}
	if currentPage < totalPage {
		htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`%d">>></a></li><li><a class="num" href="`+path+`%d">..%d</a></li>`, currentPage+1, totalPage, totalPage)
	}
	return htmlPage
}

//转换字节大小
func FormatByte(size int) string {
	fsize := float64(size)
	//字节单位
	units := [6]string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; fsize >= 1024 && i < 5; i++ {
		fsize /= 1024
	}

	num := fmt.Sprintf("%.2f", fsize)

	return string(num) + units[i]
}

//MySha512加密
func MySha512(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

//Sha256加密
func MySha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

//MD5加密函数
func MyMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//解析url链接参数，必须是键值对的形式,url链接形式:cid/4/sys/linux/wd/java/p/12，不带“.html”等后缀
func ParseURL(url string) map[string]string {
	var params map[string]string
	params = make(map[string]string, 4)
	url_slice := strings.Split(strings.Trim(url, "/"), "/")
	url_len := len(url_slice)
	if url_len > 1 {
		if url_len != (url_len/2)*2 { //处理键值非成双成对的情况
			url_len = url_len - 1
		}
		for i := 0; i < url_len; i += 2 {
			params[url_slice[i]] = url_slice[(i + 1)]
		}
	}
	return params
}

//简单加减运算,在模板上使用,如果是减法,则是第一个数减后面的数
func Calc(operate string, nums ...int) int {
	val := 0
	if operate == "+" { //加
		for _, num := range nums {
			val += num
		}
	} else { //减
		for k, num := range nums {
			if k == 0 {
				val = num
			} else {
				val -= num
			}
		}
	}
	return val
}

//生成url链接
//prefix表示前缀，如/code/xx
//subfix表示后缀，如.html
//params表示链接参数，可以是键值对或者多个值
func BuildURL(prefix string, subfix string, params ...interface{}) string {
	var (
		l    int
		url  string
		strs []string
	)
	url = "/" + strings.Trim(prefix, "/")
	l = len(params)
	if l > 0 {
		for _, v := range params {
			vv := strings.TrimSpace(fmt.Sprintf("%v", v))
			if len(vv) > 0 {
				strs = append(strs, vv)
			}
		}
	}
	str := strings.Trim(strings.Join(strs, "/"), "/")
	url = prefix + "/" + str + subfix
	return url
}

//限制数字值的范围
func NumberLimit(num, min, max int) int {
	if num < min {
		return min
	}
	if num > max {
		return max
	}
	return num
}

//整型转化成字符串
func Int2Str(num int) string {
	return strconv.Itoa(num)
}

//整型转化成字符串
func Str2Int(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//计算执行耗时
func ExecTime(s int64, ns int) string {
	spend_second := int(time.Now().Unix()) - int(s)
	spend_ns := time.Now().Nanosecond() - ns + spend_second*100000000
	spend_time := math.Ceil(float64(spend_ns)) / 100000000
	str := fmt.Sprintf("%v秒", spend_time)
	return str
}

//url请求处理
func UrlEscape(str string) string {
	return strings.TrimSpace(url.QueryEscape(strings.Replace(str, "/", " ", -1)))
}

//打印
func HCPrint(a ...interface{}) {
	//如果是开发者模式，则打印数据
	if beego.AppConfig.String("runmode") == "dev" {
		fmt.Println("")
		fmt.Println("==================")
		fmt.Println("==================")
		fmt.Println(a...)
		fmt.Println("==================")
		fmt.Println("==================")
		fmt.Println("")
	}
}

//根据ua判断访客是手机端还是PC端访问
//手机端返回true，否则返回false，默认返回false,即PC端
func IsMobile(ua string) bool {
	ua = strings.ToLower(ua)
	//常见的手机端UA判断，有即默认为手机端 pad、phone关键字替代ipad，iphone,apad等
	mobile := []string{"mobile", "phone", "android", "pad", "pod", "symbian", "wap", "smartphone", "apk", "ios"}
	for _, m := range mobile {
		if strings.Contains(ua, m) {
			return true
		}
	}
	//生僻的不常见的UA判断
	mbstr := "w3c,acs-,alav,alca,amoi,audi,avan,benq,bird,blac"
	mbstr += "blaz,brew,cell,cldc,cmd-,dang,doco,eric,hipt,inno"
	mbstr += "ipaq,java,jigs,kddi,keji,leno,lg-c,lg-d,lg-g,lge-"
	mbstr += "maui,maxo,midp,mits,mmef,mobi,mot-,moto,mwbp,nec-"
	mbstr += "newt,noki,oper,palm,pana,pant,phil,play,port,prox"
	mbstr += "qwap,sage,sams,sany,sch-,sec-,send,seri,sgh-,shar"
	mbstr += "sie-,siem,smal,smar,sony,sph-,symb,t-mo,teli,tim-"
	mbstr += "tosh,tsm-,upg1,upsi,vk-v,voda,wap-,wapa,wapi,wapp"
	mbstr += "wapr,webc,winw,winw,xda,xda-,up.browser,up.link,mmp,midp,xoom"
	slice := strings.Split(mbstr, ",")
	for _, m := range slice {
		if strings.Contains(ua, m) {
			return true
		}
	}
	return false
}

//获取指定个数的随机字符串
//size:指定的字符串个数
//kind:0,纯数字，1,小写字母，2,大写字母，3,数字+大小写字母
func RandStr(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

//获取并返回文件的md5值
func FileMd5(filepath string) (string, error) {
	file, err := os.Open(filepath)
	defer file.Close()
	if err == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		return fmt.Sprintf("%x", md5h.Sum(nil)), nil
	}
	return "", err
}

//首字母大写
func UpperFirst(str string) string {
	return strings.Replace(str, str[0:1], strings.ToUpper(str[0:1]), 1)
}

//输出默认值【模板函数】
func Default(ori, def interface{}) string {
	v := fmt.Sprintf("%v", ori)
	if len(v) == 0 {
		return fmt.Sprintf("%v", def)
	}
	return v
}

//发送邮件
func SendMail(to, subject, content string) error {
	msg := &mail.Message{
		mail.Header{
			"From":         {beego.AppConfig.String("mail_username")},
			"To":           {to},
			"Subject":      {subject},
			"Content-Type": {"text/html"},
		},
		strings.NewReader(content),
	}
	port := beego.AppConfig.DefaultInt("mail_port", 25)
	host := beego.AppConfig.String("mail_host")
	username := beego.AppConfig.String("mail_username")
	password := beego.AppConfig.String("mail_password")
	m := mailer.NewMailer(host, username, password, port)
	err := m.Send(msg)
	return err
}
