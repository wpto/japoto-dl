package muxer

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

// func NewFFMpeg(input, image, output string, album, artist, title, track *string) *FFMpeg {
// 	return &FFMpeg{
// 		Metadata: metadata{
// 			Album:  album,
// 			Artist: artist,
// 			Title:  title,
// 			Track:  track,
// 		},
// 		InputFile:  input,
// 		ImageFile:  image,
// 		OutputFile: output,
// 	}
// }

// type metadata struct {
// 	Album  *string
// 	Artist *string
// 	Title  *string
// 	Track  *string
// }

// func (m metadata) Args() []string {
// 	result := []string{}

// 	fields := []struct {
// 		key   *string
// 		label string
// 	}{
// 		{m.Album, "album"},
// 		{m.Artist, "artist"},
// 		{m.Title, "title"},
// 		{m.Track, "track"},
// 	}

// 	for _, val := range fields {
// 		if val.key != nil {
// 			result = append(result, "-metadata")
// 			result = append(result, val.label+"="+(*val.key))
// 		}
// 	}
// 	return result
// }

// type FFMpeg struct {
// 	Metadata   metadata
// 	InputFile  string
// 	ImageFile  string
// 	OutputFile string
// }

// func (ffm FFMpeg) Args() []string {
// 	result := []string{}

// 	result = append(result, "-i", ffm.InputFile)
// 	result = append(result, "-i", ffm.ImageFile)
// 	result = append(result, "-map", "0", "-map", "1:0")
// 	result = append(result, "-vc")
// 	result = append(result, "-acodec", "libmp3lame", "-q:a", "2")
// 	result = append(result, ffm.Metadata.Args()...)
// 	result = append(result, "-y") // overwrite
// 	result = append(result, ffm.OutputFile)

// 	return result
// }

func FFMpeg(args []string) error {
	fmt.Println(args)

	c := exec.Command("ffmpeg", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		fmt.Printf("ffmpeg args:\n  %v\n", args)
		fmt.Printf("ffmpeg stdout:\n  %s\n", stdout)
		fmt.Printf("ffmpeg stderr:\n  %s\n", stderr)
		return errors.Wrap(err, "ffmpeg")
	}

	return nil
}
