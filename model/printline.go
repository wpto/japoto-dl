package model

type PrintLine interface {
	SetChunk(count int)
	AddChunk()
	SetChunkCount(count int)
	SetPrefix(prefix string)
	Status(str string)
	Error(err error) error
}
