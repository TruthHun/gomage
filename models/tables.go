package models

//图片样式表
type Style struct {
	Id     int
	Sid    int    `orm:"index` //对应的系统设置的ID，即system id
	Title  string //样式标题
	Rule   string `orm:"unique"`     //样式规则名
	Width  int    `orm:"default(0)"` //图片宽度
	Height int    `orm:"default(0)"` //图片高度
	Top    int    `orm:"default(0)"` //水印顶部距离
	Left   int    `orm:"default(0)"` //水印左侧距离
	Right  int    `orm:"default(0)"` //水印右侧距离
	Bottom int    `orm:"default(0)"` //水印底部距离
	//图片裁剪处理方式Method，
	// 1、固定宽高缩放
	//11：按短边缩放，居中裁剪
	//12：按长边缩放，缩略填充
	// 2、等比例缩放
	//21:固定宽度，高度自适应
	//22:固定高度，宽度自适应
	Method            int
	Ext               string `orm:"size(30);null"` //图片后缀
	IsZoom            bool   `orm:"default(true)"` //是否允许图片放大
	Watermark         string `orm:"null"`          //水印路径
	WatermarkPosition int    `orm:"default(9)"`    //水印位置
	Status            bool   `orm:"default(true)"` //样式是否启用
	Time              int    `orm:"default(0)"`    //样式更新时间
}

// 多字段唯一键
func (s *Style) TableUnique() [][]string {
	return [][]string{
		[]string{"Sid", "Rule"},
	}
}

//管理员登录表
type Admin struct {
	Id       int
	Username string `orm:"size(32);unique"` //用户名
	Password string `orm:"size(32)"`        //密码
}

//系统设置
type System struct {
	Id        int
	Host      string `orm:"unique"`
	Sitename  string `orm:"null"`
	Referer   string `orm:"size(1024);null"`     //白名单，如果设置了白名单，则表示开启防盗链
	Prefix    string `orm:"null"`                //路径前缀
	Protected bool   `orm:"default(false)"`      //是否开启原图保护，默认不开启
	Segment   string `orm:"default(@！);size(5)"` //样式分隔符，默认“@！”，其它选项：- (中划线)_(下划线)/ (斜杠)! (感叹号)，gopic.cc/sample.jpg+分隔符+stylename
	IsCache   bool   `orm:"default(true)"`       //是否开启图片缓存，缓存的图片会放在当前目录下的cache目录
	Status    bool   `orm:"default(true)"`       //状态，true表示正常，false表示关闭
}
