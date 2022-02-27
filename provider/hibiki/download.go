package hibiki

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (ep *HibikiEpisodeMedia) Download(loader model.Loader, tasks model.Tasks) error {
	hls := tasks.AudioHLS()

	var checkObj struct {
		PlaylistUrl string `json:"playlist_url"`
	}
	err := loader.JSON("https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id="+fmt.Sprint(ep.Id), &checkObj, gopts)
	if err != nil {
		return errors.Wrap(err, "hibiki.dl.check")
	}

	playlistUrl := checkObj.PlaylistUrl
	playlistBody, err := loader.Text(playlistUrl, gopts)
	if err != nil {
		errors.Wrap(err, "hibiki.dl.playlist")
	}

	tsaudio, err := hls.Playlist(*playlistBody)
	if err != nil {
		return errors.Wrap(err, "hibiki.dl.playlist")
	}

	if len(tsaudio) > 1 {
		return errors.New("hibiki.dl.playlist: tsaudio size > 1: not implemented")
	}

	errcImg := make(chan error)
	go func(errc chan<- error) {
		url := ep.showRef.PcImageUrl
		if len(url) == 0 {
			errc <- errors.New("hibiki.dl.image: not found")
			return
		}

		imageBody, err := loader.Raw(ep.showRef.PcImageUrl, gopts)
		if err != nil {
			errc <- errors.Wrap(err, "hibiki.dl.img")
			return
		}

		file := model.NewFile("", "")
		file.SetBody(imageBody)

		err = hls.Image(file)
		if err != nil {
			errc <- errors.Wrap(err, "hibiki.dl.img")
			return
		}

		errc <- nil
	}(errcImg)

	for _, ts := range tsaudio {
		tsaudioUrl, err := ts.Url(playlistUrl)
		if err != nil {
			return errors.Wrap(err, "hibiki.dl.ts")
		}
		tsaudioBody, err := loader.Text(tsaudioUrl, gopts)
		if err != nil {
			return errors.Wrap(err, "hibiki.dl.ts")
		}
		ts.SetBodyString(tsaudioBody)

		keys, audio, err := hls.TSAudio(ts)
		if err != nil {
			return errors.Wrap(err, "hibiki.dl.ts")
		}

		for _, k := range keys {
			url, err := k.Url(tsaudioUrl)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.key")
			}

			body, err := loader.Raw(url, gopts)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.key")
			}

			k.SetBody(body)

			err = hls.Validate(k)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.key")
			}
		}

		for _, k := range audio {
			url, err := k.Url(tsaudioUrl)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.audio")
			}
			body, err := loader.Raw(url, gopts)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.audio")
			}
			k.SetBody(body)

			err = hls.Validate(k)
			if err != nil {
				return errors.Wrap(err, "hibiki.dl.audio")
			}
		}
	}

	return nil
}
