package common

import (
	"net/http"

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
