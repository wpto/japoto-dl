package types

type Entry struct {
	Date          string `json:"date"`
	Duration      int    `json:"duration"`
	DurationHuman string `json:"duration_human"`
	Filename      string `json:"filename"`
	HasImage      bool   `json:"has_image"`
	MessageId     int    `json:"message_id"`
	Performer     string `json:"performer"`
	Provider      string `json:"provider"`
	ShowId        string `json:"show_id"`
	Size          int    `json:"size"`
	SizeHuman     string `json:"size_human"`
	Title         string `json:"title"`
	URL           string `json:"url"`
}

// Performers    []Person `json:"performers"`
type Person struct {
	IsGuest   bool    `json:"is_guest"`
	Name      string  `json:"name"`
	Character *string `json:"character"`
}

type Source interface {
	Read() []Entry
	Write([]Entry)
}

type Printer interface {
	Print(db []Entry) error
}
