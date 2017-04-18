package helper

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"errors"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

//图片文件信息
type Info struct {
	Width, Height int    //图片宽高
	Size          int64  //图片文件大小
	Md5           string //图片md5
	ModTime       int    //图片修改时间戳
	Ext           string //图片后缀
}

//判断文件路径判断文件是否是符合要求的图片格式，jpeg,jpg,gif,png,bmp,tif,tiff
func IsImage(path string) bool {
	slice := strings.Split(path, ".")
	ext := strings.ToLower(strings.TrimSpace(slice[len(slice)-1]))
	exts := map[string]string{"jpeg": "jpeg", "jpg": "jpg", "gif": "gif", "png": "png", "bmp": "bmp", "tif": "tif", "tiff": "tiff"}
	_, ok := exts[ext]
	return ok
}

//获取图片文件信息
func ImageInfo(dst string) (Info, error) {
	var (
		this     Info
		fileinfo os.FileInfo
		err      error
		file     *os.File
		config   image.Config
	)
	file, err = os.Open(dst)
	defer file.Close()
	if err == nil {
		slice := strings.Split(dst, ".")
		ext := strings.ToLower(slice[len(slice)-1])
		switch ext {
		case "jpeg", "jpg":
			config, err = jpeg.DecodeConfig(file)
		case "gif":
			config, err = gif.DecodeConfig(file)
		case "png":
			config, err = png.DecodeConfig(file)
		case "bmp":
			config, err = bmp.DecodeConfig(file)
		case "tif", "tiff":
			config, err = tiff.DecodeConfig(file)
		default:
			err = errors.New("Not an image format")
		}
		if err == nil {
			this.Width = config.Width
			this.Height = config.Height
			this.Ext = ext
			fileinfo, err = os.Stat(dst)
			if err == nil {
				this.Size = fileinfo.Size()
				this.ModTime = int(fileinfo.ModTime().Unix())
				md5h := md5.New()
				io.Copy(md5h, file)
				this.Md5 = fmt.Sprintf("%x", md5h.Sum(nil))
			}
		}
	}
	return this, err
}
