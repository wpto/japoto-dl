package cmd

import "fmt"

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
}

func (pl *PrintLine) appendSpace(str string) string {
	for i := 0; i < (pl.lastLen - len(str)); i++ {
		str += " "
	}
	return str
}

func (pl *PrintLine) Error(err error) error {
	res := fmt.Sprintf("\r%s: %s\n", pl.prefix, err.Error())
	fmt.Print(pl.appendSpace(res))
	return err
}

func (pl *PrintLine) Print(str string) {
	fmt.Printf("\r%s\n", pl.appendSpace(str))
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

	newLen := len(str)
	str = pl.appendSpace(str)
	pl.lastLen = newLen

	fmt.Printf("\r%s", str)
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
