package onsen

import (
	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/types"
	"github.com/pkg/errors"
)

type Onsen struct {
	loader types.Loader
	tasks  *tasks.Tasks
}

func NewOnsen(loader types.Loader, tasks *tasks.Tasks) *Onsen {
	return &Onsen{loader, tasks}
}

var gopts *grequests.RequestOptions = &grequests.RequestOptions{
	Headers: map[string]string{
		"Referer": "https://www.onsen.ag/",
	},
}

func (p *Onsen) Download(playlistUrl string) error {

	playlistBody, err := p.loader.Text(playlistUrl, gopts)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plget")
	}

	tsaudio, err := p.tasks.AudioHLS.Playlist(*playlistBody)
	if err != nil {
		return errors.Wrap(err, "onsen.dl.plparse")
	}

	for _, ts := range tsaudio {
		tsaudioUrl, err := ts.Url(playlistUrl)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsurl")
		}
		tsaudioBody, err := p.loader.Text(tsaudioUrl, gopts)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsget")
		}

		ts.SetBodyString(tsaudioBody)

		keys, audio, err := p.tasks.AudioHLS.TSAudio(ts)
		if err != nil {
			return errors.Wrap(err, "onsen.dl.tsparse")
		}

		for _, key := range keys {
			url, err := key.Url(tsaudioUrl)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyurl")
			}

			body, err := p.loader.Raw(url, gopts)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyget")
			}

			key.SetBody(body)

			err = p.tasks.AudioHLS.Validate(key)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyuse")
			}
		}

		for _, au := range audio {
			url, err := au.Url(tsaudioUrl)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyurl")
			}

			body, err := p.loader.Raw(url, gopts)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyget")
			}

			au.SetBody(body)

			err = p.tasks.AudioHLS.Validate(au)
			if err != nil {
				return errors.Wrap(err, "onsen.dl.keyuse")
			}
		}
	}

	return nil
}
