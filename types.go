package imgproxygo

import "fmt"

type imgProxy struct {
	config Config
}

type imgproxyURL struct {
	ext   FORMAT_TYPE
	paths []string
	*imgProxy
}

type gravity struct {
	genre GRAVITY_TYPE
	x     any
	y     any
}

func (g *gravity) String() string {
	if g.genre == GRAVITY_TYPE_FOCUS {
		if g.x != nil && g.y != nil {
			return fmt.Sprintf("%s:%s:%s", g.genre, cleanFloat(*g.x.(*float64)), cleanFloat(*g.y.(*float64)))
		}
	} else if g.genre != GRAVITY_TYPE_SMART {
		if g.x != nil && g.y != nil {
			return fmt.Sprintf("%s:%d:%d", g.genre, *g.x.(*int), *g.y.(*int))
		}
	}
	return string(g.genre)
}

type RESIZE_TYPE string

const (
	RESIZE_TYPE_FIT       RESIZE_TYPE = "fit"
	RESIZE_TYPE_FILL      RESIZE_TYPE = "fill"
	RESIZE_TYPE_FILL_DOWN RESIZE_TYPE = "fill-down"
	RESIZE_TYPE_FORCE     RESIZE_TYPE = "force"
	RESIZE_TYPE_AUTO      RESIZE_TYPE = "auto"
)

type GRAVITY_TYPE string

const (
	GRAVITY_TYPE_NO    GRAVITY_TYPE = "no"   //上对齐
	GRAVITY_TYPE_SO    GRAVITY_TYPE = "so"   //下对齐
	GRAVITY_TYPE_EA    GRAVITY_TYPE = "ea"   //右对齐
	GRAVITY_TYPE_WE    GRAVITY_TYPE = "we"   //左对齐
	GRAVITY_TYPE_NOEA  GRAVITY_TYPE = "noea" //右上
	GRAVITY_TYPE_NOWE  GRAVITY_TYPE = "nowe" //左上
	GRAVITY_TYPE_SOEA  GRAVITY_TYPE = "soea" //右下
	GRAVITY_TYPE_SOWE  GRAVITY_TYPE = "sowe" //左下
	GRAVITY_TYPE_CE    GRAVITY_TYPE = "ce"   //中央
	GRAVITY_TYPE_SMART GRAVITY_TYPE = "sm"   //特殊，智能判断
	GRAVITY_TYPE_FOCUS GRAVITY_TYPE = "fp"   //特殊，使用0~1的浮点数，表示从右到左或者从上到下的中央点
)

type FORMAT_TYPE string

const (
	FORMAT_TYPE_PNG  FORMAT_TYPE = "png"
	FORMAT_TYPE_JPEG FORMAT_TYPE = "jpg"
	FORMAT_TYPE_WEBP FORMAT_TYPE = "webp"
	FORMAT_TYPE_AVIF FORMAT_TYPE = "avif"
	FORMAT_TYPE_GIF  FORMAT_TYPE = "gif"
	FORMAT_TYPE_ICO  FORMAT_TYPE = "ico"
	FORMAT_TYPE_BMP  FORMAT_TYPE = "bmp"
	FORMAT_TYPE_TIFF FORMAT_TYPE = "tiff"
)
