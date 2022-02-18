package onsen

import "github.com/pgeowng/japoto-dl/types"

type OnsenInfo struct {
	loader types.Loader
}

func NewOnsenInfo(loader types.Loader) *OnsenInfo {
	return &OnsenInfo{loader}
}
