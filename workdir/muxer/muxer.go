package muxer

import "fmt"

type MuxerHLS interface {
	Mux(inputPath string) error
}

type MuxerHLSImpl struct {
	outputPath string
	imagePath  *string
	tags       map[string]string
}

func NewMuxerHLS(outputPath string, imagePath *string, tags map[string]string) MuxerHLS {
	return &MuxerHLSImpl{outputPath, imagePath, tags}
}

func (m *MuxerHLSImpl) Mux(inputPath string) error {
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
