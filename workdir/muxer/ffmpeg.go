package muxer

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type FFMpegHLS struct {
	outputPath string
	imagePath  *string
	tags       map[string]string
}

func NewFFMpegHLS(outputPath string, imagePath *string, tags map[string]string) *FFMpegHLS {
	return &FFMpegHLS{outputPath, imagePath, tags}
}

func (m *FFMpegHLS) Mux(inputPath string) error {
	result := []string{}

	result = append(result, "-allowed_extensions", "ALL")
	result = append(result, "-i", inputPath)
	if m.imagePath == nil {
		result = append(result, "-vn")
	} else {
		result = append(result, "-i", *m.imagePath)
		result = append(result, "-map", "0", "-map", "1:0")
		result = append(result, "-vc")
	}
	result = append(result, "-acodec", "libmp3lame", "-q:a", "2")

	for tag, val := range m.tags {
		result = append(result, "-metadata")
		result = append(result, fmt.Sprintf("%s=%s", tag, val))
	}

	result = append(result, "-y") // overwrite
	result = append(result, m.outputPath)

	return FFMpeg(result)
}

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
