package provider

type Onsen struct{}

func NewOnsen() *Onsen {
	return &Onsen{}
}

func (o *Onsen) Label() string {
	return "onsen"
}
