package audiohls

import (
	"fmt"
	"regexp"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pkg/errors"
)

func (a *AudioHLSImpl) TSAudio(tsaudio model.File) (keys []model.File, audio []model.File, err error) {
	keyRE := regexp.MustCompile(`(?m)(#EXT-X-KEY:METHOD=AES-128,URI=")([^"]*)([^\n]*)`)
	keys = make([]model.File, 0)
	idx := 0

	myTSAudioText := keyRE.ReplaceAllStringFunc(tsaudio.BodyString(), func(line string) string {
		match := keyRE.FindStringSubmatch(line)
		if match == nil {
			panic("replace match failed")
		}

		link := match[2]
		name := fmt.Sprintf("key_%d", idx)
		idx += 1

		keys = append(keys, model.NewFile(link, name))

		return match[1] + name + match[3]
	})

	audioRE := regexp.MustCompile(`(?m)^([^#]+?)$`)
	audio = make([]model.File, 0)
	idx = 0

	myTSAudioText = audioRE.ReplaceAllStringFunc(myTSAudioText, func(link string) string {
		name := fmt.Sprintf("audio_%d.ts", idx)
		idx += 1
		audio = append(audio, model.NewFile(link, name))
		return name
	})

	if len(audio) == 0 {
		fmt.Println("ahls.tsaudio: links not found")
		fmt.Println(myTSAudioText)
	}

	err = a.workdir.Save(tsaudio.Name(), myTSAudioText)
	if err != nil {
		return nil, nil, errors.Wrap(err, "ahls.tsaudio")
	}

	return keys, audio, nil
}
