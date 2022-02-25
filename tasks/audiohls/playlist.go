package audiohls

import (
	"fmt"
	"log"
	"regexp"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (a *AudioHLSImpl) Playlist(playlistText string) (tsaudio []model.File, err error) {
	re := regexp.MustCompile(`(?m)^([^#]+?)$`)
	tsaudio = make([]model.File, 0)
	idx := 0

	myPlaylistText := re.ReplaceAllStringFunc(playlistText, func(link string) string {
		name := fmt.Sprintf("tsaudio_%d.m3u8", idx)
		idx += 1
		tsaudio = append(tsaudio, model.NewFile(link, name))
		return name
	})

	if idx == 0 {
		fmt.Println(playlistText)
		return nil, errors.New("tsaudio links not found")
	}

	if idx > 1 {
		log.Printf("ahls.playlist: has more then 1 tsaudio: %v", tsaudio)
		fmt.Println(playlistText)
	}

	err = a.workdir.SaveNamed("playlist", myPlaylistText)
	if err != nil {
		return nil, errors.Wrap(err, "ahls.playlist")
	}

	return tsaudio, nil
}
