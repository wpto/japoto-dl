package usecase

import "github.com/pgeowng/japoto-dl/provider/onsen"

type DownloadUsecase struct {
	onsen onsen.OnsenUsecase
}

func NewDownload() *DownloadUsecase {
	return &DownloadUsecase{}
}

func RunOnsen() {

}
