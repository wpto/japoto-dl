package audiohls

func (a *AudioHLSImpl) CheckAlreadyLoaded(filename string) bool {
	loaded := a.workdir.AlreadyLoadedChunks()
	return loaded[filename]
}
