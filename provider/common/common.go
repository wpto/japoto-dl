package common

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pgeowng/japoto-dl/model"
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

func LoadPlaylist(playlistUrl string, gopts *model.LoaderOpts, loader model.Loader, hls model.AudioHLS) (result []model.File, err error) {
	playlistBody, err := loader.Text(playlistUrl, gopts)
	if err != nil {
		err = errors.Wrap(err, "load playlist")
		return
	}

	result, err = hls.Playlist(*playlistBody)
	if err != nil {
		err = errors.Wrap(err, "save playlist")
		return
	}

	return
}

func LoadTSAudio(playlistUrl string, gopts *model.LoaderOpts, ts model.File, loader model.Loader, hls model.AudioHLS) (keys []model.File, audio []model.File, tsaudioUrl string, err error) {
	tsaudioUrl, err = ts.Url(playlistUrl)
	if err != nil {
		err = errors.Wrap(err, "onsen.dl.tsurl")
		return
	}

	tsaudioBody, err := loader.Text(tsaudioUrl, gopts)
	if err != nil {
		err = errors.Wrap(err, "onsen.dl.tsget")
		return
	}

	ts.SetBodyString(tsaudioBody)

	keys, audio, err = hls.TSAudio(ts)
	if err != nil {
		err = errors.Wrap(err, "onsen.dl.tsparse")
		return
	}

	return
}

func FilterChunks(src []model.File, hls model.AudioHLS) []model.File {
	result := []model.File{}
	for _, file := range src {
		if !hls.CheckAlreadyLoaded(file.Name()) {
			result = append(result, file)
		}
	}
	return result
}
