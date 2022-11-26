package model

type Tags struct {
	Album  *string
	Artist *string
	Title  *string
	Track  *string
}

func NewTags(ep Episode) *Tags {

	return &Tags{}
}
