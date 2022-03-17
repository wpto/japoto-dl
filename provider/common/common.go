package common

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ext map[string]string = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
}

func GuessContentType(file []byte) string {
	buf := make([]byte, 512)
	_ = copy(buf, file)
	contentType := http.DetectContentType(buf)
	result, ok := ext[contentType]
	if !ok {
		panic(errors.Errorf("guess type: not found %s", contentType))
	}
	return result
}

func EncodeIdx(f int, nums ...int) string {
	first := strconv.FormatInt(int64(f), 35)
	if len(nums) == 0 {
		return first
	}

	rest := make([]string, 0, len(nums)+1)
	rest = append(rest, first)
	for _, v := range nums {
		rest = append(rest, strconv.FormatInt(int64(v), 35))
	}

	return strings.Join(rest, "z")
}
