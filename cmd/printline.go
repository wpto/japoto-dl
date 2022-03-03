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

func (pl *PrintLine) Error(err error) error {
	fmt.Printf("\r%s: %s\n", pl.prefix, err.Error())
	// pl.printStatus()
	return err
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
	if pl.lastLen > newLen {
		for i := 0; i < (pl.lastLen - newLen); i++ {
			str += " "
		}
	}

	pl.lastLen = newLen

	fmt.Printf("\r%s", str)
}
