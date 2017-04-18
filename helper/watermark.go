package helper

import (
	"image"
	"strings"

	"github.com/disintegration/imaging"
)

//给图片添加水印
const (
	WATERMARK_POSITION_NORTHWEST int = 1        //左上角水印
	WATERMARK_POSITION_NORTH     int = iota + 1 //上居中水印
	WATERMARK_POSITION_NORTHEAST                //右上角水印
	WATERMARK_POSITION_WEST                     //左居中水印
	WATERMARK_POSITION_CENTER                   //居中水印
	WATERMARK_POSITION_EAST                     //右居中水印
	WATERMARK_POSITION_SOUTHWEST                //左下角水印
	WATERMARK_POSITION_SOUTH                    //下居中水印
	WATERMARK_POSITION_SOUTHEAST                //右下角水印
)

//添加水印
//img:需要添加水印的图片对象
//watermark:水印的图片路径
//position:水印位置，共有9个
//top, right, bottom, left：水印距离
func AddWatermark(img image.Image, watermark string, position, top, right, bottom, left int) (image.Image, error) {
	var (
		err   error
		point image.Point
	)
	watermark = strings.TrimLeft(watermark, "./")
	img_watermark, err := imaging.Open(watermark)
	if err == nil {
		//水印图片的宽高
		MaxWater, MinWater := img_watermark.Bounds().Max, img_watermark.Bounds().Min
		WidthWater, HeightWater := MaxWater.X-MinWater.X, MaxWater.Y-MinWater.Y
		//目标图片的宽高
		Max, Min := img.Bounds().Max, img.Bounds().Min
		Width, Height := Max.X-Min.X, Max.Y-Min.Y
		//+top,+left,-right,-bottom
		switch position {
		case WATERMARK_POSITION_NORTHWEST:
			point = image.Point{0 + left, 0 + top}
		case WATERMARK_POSITION_NORTH:
			point = image.Point{(Width - WidthWater) / 2, 0 + top}
		case WATERMARK_POSITION_NORTHEAST:
			point = image.Point{Width - WidthWater - right, 0 + top}
		case WATERMARK_POSITION_WEST:
			point = image.Point{0 + left, (Height - HeightWater) / 2}
		case WATERMARK_POSITION_CENTER:
			point = image.Point{(Width - WidthWater) / 2, (Height - HeightWater) / 2}
		case WATERMARK_POSITION_EAST:
			point = image.Point{Width - WidthWater - right, (Height - HeightWater) / 2}
		case WATERMARK_POSITION_SOUTHWEST:
			point = image.Point{0 + left, Height - HeightWater - bottom}
		case WATERMARK_POSITION_SOUTH:
			point = image.Point{(Width - WidthWater) / 2, Height - HeightWater - bottom}
		case WATERMARK_POSITION_SOUTHEAST:
			point = image.Point{Width - WidthWater - right, Height - HeightWater - bottom}
		}
		img = imaging.Paste(img, img_watermark, point)
	}
	return img, err
}
