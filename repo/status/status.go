package status

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type PrintLine struct {
	chunk      int
	chunkCount int

	loaded      int
	loadedCount int

	muxed      int
	muxedCount int

	lastLen int

	status string
	prefix string

	w io.Writer
}

func New(writer io.Writer) *PrintLine {
	return &PrintLine{w: writer}
}

func appendSpace(str string, targetSize int) string {
	for i := targetSize - len(str); i > 0; i-- {
		str += " "
	}
	return str
}

func (pl *PrintLine) Error(err error) error {
	res := fmt.Sprintf("%s: %s", pl.prefix, err.Error())
	pl.Println(res)
	return err
}

func (pl *PrintLine) Print(str string) {
	newLen := len(str)
	str = appendSpace(str, pl.lastLen)
	fmt.Fprintf(pl.w, "\r%s", str)
	pl.lastLen = newLen
}

func (pl *PrintLine) Println(str string) {
	str = appendSpace(str, pl.lastLen)
	fmt.Fprintf(pl.w, "\r%s\n", str)
	pl.lastLen = 0
}

func (pl *PrintLine) AddLoaded() {
	pl.loaded++
	pl.printStatus()
}

func (pl *PrintLine) AddMuxed() {
	pl.muxed++
	pl.printStatus()
}

func (pl *PrintLine) AddMuxedCount() {
	pl.muxedCount++
	pl.printStatus()
}
func (pl *PrintLine) AddLoadedCount() {
	pl.loadedCount++
	pl.printStatus()
}

func (pl *PrintLine) SetChunk(count int) {
	pl.chunk = count
	pl.printStatus()
}

func (pl *PrintLine) AddChunk() {
	pl.chunk++
	pl.printStatus()
}

func (pl *PrintLine) SetChunkCount(count int) {
	pl.chunkCount = count
	pl.printStatus()
}

func (pl *PrintLine) SetPrefix(prefix string) {
	pl.prefix = prefix
	pl.printStatus()
}

func (pl *PrintLine) Status(status string) {
	pl.status = status
	pl.printStatus()
}

func (pl *PrintLine) printStatus() {
	str := fmt.Sprintf("[%d/%d %d/%d %d/%d] %s: %s",
		pl.muxed, pl.muxedCount,
		pl.loaded, pl.loadedCount,
		pl.chunk, pl.chunkCount,
		pl.prefix, pl.status,
	)

	pl.Print(str)
}

func (pl *PrintLine) SetEp(provider, epidx string) {
	pl.SetPrefix(fmt.Sprintf("%s/%s", provider, epidx))
	pl.SetChunk(0)
	pl.printStatus()
}

func (pl *PrintLine) Inc(step int) {
	pl.SetChunk(pl.chunk + step)
	pl.printStatus()
}

func (pl *PrintLine) Total(total int) {
	// TODO: total already set
	pl.SetChunkCount(total)
	pl.printStatus()
}

type ErrorPrintLine struct{}

func (e *ErrorPrintLine) SetChunk(count int)      {}
func (e *ErrorPrintLine) AddChunk()               {}
func (e *ErrorPrintLine) SetChunkCount(count int) {}
func (e *ErrorPrintLine) SetPrefix(prefix string) {}
func (e *ErrorPrintLine) Status(str string)       {}
func (e *ErrorPrintLine) Error(err error) error {
	fmt.Printf("%s\n", err.Error())
	return err
}

type Metric interface {
	Set(label string, value float32)
	Add(label string, value float32)
	Inc(label string)
	WithLabel(label string, value string) Metric
}

type MetricImpl struct {
}

func NewMetric() Metric {
	return &MetricImpl{}
}

func (m *MetricImpl) Set(label string, value float32) {
	log.Debug().Str("metric", label).Float32("set", value).Msg("set metric")
}

func (m *MetricImpl) Add(label string, value float32) {
	log.Debug().Str("metric", label).Float32("add", value).Msg("add metric")
}

func (m *MetricImpl) Inc(label string) {
	log.Debug().Str("metric", label).Float32("inc", 1.0).Msg("inc metric")
}

func (m *MetricImpl) WithLabel(label string, value string) Metric {
	log.Debug().Str("label", label).Str("value", value).Msg("set label")
	return m
}
