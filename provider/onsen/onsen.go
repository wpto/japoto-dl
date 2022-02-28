package onsen

type Onsen struct{}

func NewOnsen() *Onsen {
	return &Onsen{}
}

func (o *Onsen) Label() string {
	return "onsen"
}
