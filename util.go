package imgproxygo

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type number interface {
	int | uint | int64 | float32 | float64
}

func minValuePtr[T number](value T, args ...T) *T {
	var (
		min T
	)
	if len(args) == 0 {
		min = 0
	} else {
		min = args[0]
	}
	if value < min {
		return nil
	}
	return &value
}

func rangeValuePtr[T number](value, min, max T) *T {
	if value >= min && value <= max {
		return &value
	}
	return nil
}

func Gravity(gravityType GRAVITY_TYPE, args ...any) (g gravity) {
	g.genre = gravityType
	if len(args) == 2 {
		if g.genre == GRAVITY_TYPE_FOCUS {
			if x, ok := args[0].(float64); ok {
				if y, ok := args[1].(float64); ok {
					g.x = rangeValuePtr(x, 0, 1)
					g.y = rangeValuePtr(y, 0, 1)
				}
			}
		} else {
			if x, ok := args[0].(int); ok {
				if y, ok := args[1].(int); ok {
					g.x = minValuePtr(x)
					g.y = minValuePtr(y)
				}
			}
		}
	}
	return g
}

func buildImgUrl(imgUrl string) string {
	url := base64.RawURLEncoding.EncodeToString([]byte(imgUrl))
	urlchunk := chunkString(url, 16)
	return strings.Join(urlchunk, "/")
}

func chunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
func boolToStr(v bool) string {
	if v {
		return "1"
	}
	return ""
}

func cleanFloat[T float32 | float64](val T) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", val), "0"), ".")
}
