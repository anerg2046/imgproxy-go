package imgproxygo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"
)

func N(config Config) *imgProxy {
	config.BaseUrl = strings.TrimRight(config.BaseUrl, "/")
	return &imgProxy{
		config: config,
	}
}

func (i *imgProxy) Builder() *imgproxyURL {
	return &imgproxyURL{
		imgProxy: i,
		ext:      "jpg",
	}
}

// ===============================
// Set Simple Options
// ===============================

// When set to `1`, `t` or `true`, imgproxy will automatically rotate images based on the EXIF Orientation parameter
// (if available in the image meta data). The orientation tag will be removed from the image in all cases.
// Normally this is controlled by the IMGPROXY_AUTO_ROTATE configuration but this procesing option allows the
// configuration to be set for each request.
//
// 当设置为`1`, `t`或`true`时，imgproxy 将根据 `EXIF` 方向参数（如果在图像元数据中可用）自动旋转图像。
// 在所有情况下，都会从图像中删除方向标签。通常这由IMGPROXY_AUTO_ROTATE配置控制，但此处理选项允许为每个请求设置配置。
func (i *imgproxyURL) AutoRotate(auto_rotate bool) *imgproxyURL {
	if auto_rotate {
		i.paths = append(i.paths, fmt.Sprintf("ar:%d", boolToInt(auto_rotate)))
	}
	return i
}

// When set, imgproxy will apply a gaussian blur filter to the resulting image.
// The value of `sigma` defines the size of the mask imgproxy will use.
//
// 设置后，imgproxy将对生成的图像应用高斯模糊。`sigma`的值越大就越模糊。
//
// 	Default: disabled
func (i *imgproxyURL) Blur(sigma float32) *imgproxyURL {
	if sigma > 0 {
		i.paths = append(i.paths, fmt.Sprintf("bl:%s", cleanFloat(sigma)))
	}
	return i
}

// Cache buster doesn’t affect image processing but its changing allows for bypassing the CDN,
// proxy server and browser cache. Useful when you have changed some things that are not reflected in the URL,
// like image quality settings, presets, or watermark data.
//
// It’s highly recommended to prefer the cachebuster option over a URL query string because that option can be properly signed.
//
// 简单点说，就是在路径里增加一个字符串，这样你就不用去刷新cdn的缓存了，因为会变成一个新地址
// 另外，这个参数支持url签名，官方建议图片有变动的时候增加这个，而不是使用url的查询字符串，因为cdn可能忽略查询字符串
//
// 	Default: empty
func (i *imgproxyURL) Cachebuster(cb string) *imgproxyURL {
	if cb != "" {
		i.paths = append(i.paths, fmt.Sprintf("cb:%s", cb))
	}
	return i
}

// When set, imgproxy will multiply the image dimensions according to this factor for HiDPI (Retina) devices.
// The value must be greater than 0.
//
// 📝Note: dpr also sets the Content-DPR header in the response so the browser can correctly render the image.
//
// 主要是针对视网膜屏这类的，设置一个比例，就是苹果App里用的2倍图，3倍图的概念
//
// 	Default: 1
func (i *imgproxyURL) Dpr(dpr float32) *imgproxyURL {
	if dpr > 0 {
		i.paths = append(i.paths, fmt.Sprintf("dpr:%s", cleanFloat(dpr)))
	}
	return i
}

// When set to 1, t or true and the source image has an embedded thumbnail, imgproxy will always use the embedded
// thumbnail instead of the main image. Currently, only thumbnails embedded in heic and avif are supported.
// This is normally controlled by the IMGPROXY_ENFORCE_THUMBNAIL configuration but this procesing option allows
// the configuration to be set for each request.
//
// 当设置为1、t或true且源图像有嵌入的缩略图时，imgproxy将始终使用嵌入的缩略图而不是主图像。
// 目前只支持嵌入heic和avif中的缩略图。这通常由IMGPROXY_ENFORCE_THUMBNAIL配置控制，
// 但是这个处理选项允许为每个请求设置配置。
func (i *imgproxyURL) EnforceThumbnail(eth bool) *imgproxyURL {
	if eth {
		i.paths = append(i.paths, fmt.Sprintf("eth:%d", boolToInt(eth)))
	}
	return i
}

// When set to 1, t or true, imgproxy will enlarge the image if it is smaller than the given size.
//
// 如果源图尺寸比要调整的小，当此参数被设置为true，就会把源图进行放大操作，结果就是图片会模糊一些
//
//  Default: false
func (i *imgproxyURL) EnLarge(el bool) *imgproxyURL {
	if el {
		i.paths = append(i.paths, fmt.Sprintf("el:%d", boolToInt(el)))
	}
	return i
}

// When set, imgproxy will check the provided unix timestamp and return 404 when expired.
//
// 设置一个unix timstamp时间戳，到了该时间，再访问图片就会返回404错误
//
//  Default: empty
func (i *imgproxyURL) Expires(exp int64) *imgproxyURL {
	if exp > time.Now().Unix() {
		i.paths = append(i.paths, fmt.Sprintf("exp:%d", exp))
	}
	return i
}

// Defines a filename for the Content-Disposition header. When not specified, imgproxy will get the filename from the source url.
//
// 当使用右键保存或者header发送下载信息时候，需要设置的图片文件名，不指定的话就取源图的文件名
//
// Default: empty
func (i *imgproxyURL) Filename(fn string) *imgproxyURL {
	if fn != "" {
		i.paths = append(i.paths, fmt.Sprintf("fn:%s", fn))
	}
	return i
}

// Specifies the resulting image format. Alias for the extension part of the URL.
//
// 指定要返回的图片文件格式
//
// Default: jpg
func (i *imgproxyURL) Format(ext FORMAT_TYPE) *imgproxyURL {
	i.ext = ext
	return i
}

// Defines the height of the resulting image. When set to 0, imgproxy will calculate resulting height using the defined width and source aspect ratio.
// When set to 0 and resizing type is force, imgproxy will keep the original height.
//
// 设置图像的高度。当设置为 0 时，imgproxy 将使用设置的宽度和源纵横比计算结果高度。
// 当设置为 0 并且调整大小类型为 force 时，imgproxy 将保持原来的高度。
//
// Default: 0
func (i *imgproxyURL) Height(h uint) *imgproxyURL {
	if h > 0 {
		i.paths = append(i.paths, fmt.Sprintf("h:%d", h))
	}
	return i
}

func (i *imgproxyURL) KeepCopyright(kcr bool) *imgproxyURL {
	if kcr {
		i.paths = append(i.paths, fmt.Sprintf("kcr:%d", boolToInt(kcr)))
	}
	return i
}

func (i *imgproxyURL) MaxBytes(mb uint) *imgproxyURL {
	if mb > 0 {
		i.paths = append(i.paths, fmt.Sprintf("mb:%d", mb))
	}
	return i
}

func (i *imgproxyURL) MinHeight(mh uint) *imgproxyURL {
	if mh > 0 {
		i.paths = append(i.paths, fmt.Sprintf("mh:%d", mh))
	}
	return i
}

func (i *imgproxyURL) MinWidth(mw uint) *imgproxyURL {
	if mw > 0 {
		i.paths = append(i.paths, fmt.Sprintf("mw:%d", mw))
	}
	return i
}
func (i *imgproxyURL) Pixelate(pix uint) *imgproxyURL {
	if pix > 0 {
		i.paths = append(i.paths, fmt.Sprintf("pix:%d", pix))
	}
	return i
}
func (i *imgproxyURL) Preset(pr string) *imgproxyURL {
	if pr != "" {
		i.paths = append(i.paths, fmt.Sprintf("pr:%s", pr))
	}
	return i
}

func (i *imgproxyURL) Quality(q float32) *imgproxyURL {
	if q > 0 && q <= 100 {
		i.paths = append(i.paths, fmt.Sprintf("q:%s", cleanFloat(q)))
	}
	return i
}

func (i *imgproxyURL) ResizingType(rt RESIZE_TYPE) *imgproxyURL {
	i.paths = append(i.paths, fmt.Sprintf("rt:%s", rt))
	return i
}

func (i *imgproxyURL) ReturnAttachment(att bool) *imgproxyURL {
	if att {
		i.paths = append(i.paths, fmt.Sprintf("att:%d", boolToInt(att)))
	}
	return i
}
func (i *imgproxyURL) Rotate(rot uint) *imgproxyURL {
	if rot >= 360 {
		rot -= 360
	}
	if rot%90 != 0 {
		rot = 0
	}
	if rot > 0 {
		i.paths = append(i.paths, fmt.Sprintf("rot:%d", rot))
	}
	return i
}

func (i *imgproxyURL) Sharpen(sh float32) *imgproxyURL {
	if sh > 0 {
		i.paths = append(i.paths, fmt.Sprintf("sh:%s", cleanFloat(sh)))
	}
	return i
}

func (i *imgproxyURL) StripMetadata(sm bool) *imgproxyURL {
	if sm {
		i.paths = append(i.paths, fmt.Sprintf("sm:%d", boolToInt(sm)))
	}
	return i
}

func (i *imgproxyURL) Width(w uint) *imgproxyURL {
	if w > 0 {
		i.paths = append(i.paths, fmt.Sprintf("w:%d", w))
	}
	return i
}

// ===============================
// Set Muti Params Options
// ===============================

func (i *imgproxyURL) Background(args ...any) *imgproxyURL {
	if len(args) == 1 {
		if hex, ok := args[0].(string); ok {
			i.paths = append(i.paths, fmt.Sprintf("bg:%s", hex))
		}
	}
	if len(args) == 3 {
		if r, ok := args[0].(int); ok {
			if g, ok := args[1].(int); ok {
				if b, ok := args[2].(int); ok {
					if r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255 {
						i.paths = append(i.paths, fmt.Sprintf("bg:%d:%d:%d", r, g, b))
					}
				}
			}
		}
	}
	return i
}

func (i *imgproxyURL) Crop(width, height float32, args ...gravity) *imgproxyURL {
	if width > 0 && height > 0 {
		path := fmt.Sprintf("c:%s:%s", cleanFloat(width), cleanFloat(height))
		if len(args) == 1 {
			path = fmt.Sprintf("%s:%s", path, args[0].String())
		}
		i.paths = append(i.paths, path)
	}
	return i
}

func (i *imgproxyURL) Extend(ex bool, args ...gravity) *imgproxyURL {
	if ex {
		path := fmt.Sprintf("ex:%d", boolToInt(ex))
		if len(args) == 1 {
			path = fmt.Sprintf("%s:%s", path, args[0].String())
		}
		i.paths = append(i.paths, path)
	}
	return i
}

func (i *imgproxyURL) Gravity(g gravity) *imgproxyURL {
	i.paths = append(i.paths, fmt.Sprintf("g:%s", g.String()))
	return i
}

func (i *imgproxyURL) Padding(args ...int) *imgproxyURL {
	if len(args) > 0 {
		if len(args) == 1 {
			i.paths = append(i.paths, fmt.Sprintf("pd:%d", args[0]))
		} else if len(args) == 2 {
			i.paths = append(i.paths, fmt.Sprintf("pd:%d:%d", args[0], args[1]))
		} else if len(args) == 3 {
			i.paths = append(i.paths, fmt.Sprintf("pd:%d:%d:%d", args[0], args[1], args[2]))
		} else if len(args) == 4 {
			i.paths = append(i.paths, fmt.Sprintf("pd:%d:%d:%d:%d", args[0], args[1], args[2], args[3]))
		}
	}
	return i
}

func (i *imgproxyURL) Trim(threshold float32, args ...any) *imgproxyURL {
	if threshold > 0 {
		path := fmt.Sprintf("t:%s", cleanFloat(threshold))
		if len(args) > 0 {
			if color, ok := args[0].(string); ok {
				if _, err := hex.DecodeString(args[0].(string)); err == nil && len(args[0].(string)) == 6 {
					path = fmt.Sprintf("%s:%s", path, color)
				}
			}
			if len(args) == 2 {
				if equal_hor, ok := args[1].(bool); ok && equal_hor {
					path = fmt.Sprintf("%s:%d", path, boolToInt(equal_hor))
				}
			}

			if len(args) == 3 {
				if equal_hor, ok := args[1].(bool); ok {
					if equal_ver, ok := args[2].(bool); ok {
						path = fmt.Sprintf("%s:%s:%s", path, boolToStr(equal_hor), boolToStr(equal_ver))
					}
				}

			}

		}
		i.paths = append(i.paths, strings.TrimRight(path, ":"))
	}
	return i
}

func (i *imgproxyURL) Zoom(args ...float32) *imgproxyURL {
	if len(args) == 1 {
		i.paths = append(i.paths, fmt.Sprintf("z:%s", cleanFloat(args[0])))
	} else if len(args) == 2 {
		i.paths = append(i.paths, fmt.Sprintf("z:%s:%s", cleanFloat(args[0]), cleanFloat(args[1])))
	}
	return i
}

func (i *imgproxyURL) Gen(imgUrl string) string {
	if i.config.Encode {
		imgUrl = buildImgUrl(imgUrl)
	}
	path := fmt.Sprintf("/%s/%s.%s", strings.Join(i.paths, "/"), imgUrl, i.ext)
	sig := i.signature(path)
	return fmt.Sprintf("%s/%s%s", i.config.BaseUrl, sig, path)
}

func (i *imgproxyURL) signature(value string) string {
	if i.config.Key == "" || i.config.Salt == "" {
		return "unsafe"
	}
	signatureSize := 32
	if i.config.SignatureSize > 0 {
		signatureSize = i.config.SignatureSize
	}
	var keyBin, saltBin []byte
	var err error

	if keyBin, err = hex.DecodeString(i.config.Key); err != nil {
		log.Fatal("Key expected to be hex-encoded string")
	}

	if saltBin, err = hex.DecodeString(i.config.Salt); err != nil {
		log.Fatal("Salt expected to be hex-encoded string")
	}
	mac := hmac.New(sha256.New, keyBin)
	mac.Write(saltBin)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)[:signatureSize])
}
