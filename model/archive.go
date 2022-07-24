package model

type ArchiveItem struct {
	HistoryKey  string                  `json:"history_key"`
	Description *ArchiveItemDescription `json:"desc,omitempty"`
	Meta        *ArchiveItemMeta        `json:"meta,omitempty"`
	Chan        *ArchiveItemChan        `json:"chan,omitempty"`
}

type ArchiveItemDescription struct {
	Date      string   `json:"date"`
	Source    string   `json:"source"`
	ShowId    string   `json:"show_id"`
	ShowTitle string   `json:"show_title"`
	EpTitle   string   `json:"ep_title"`
	Artists   []string `json:"artists"`
}

type ArchiveItemMeta struct {
	Filename string `json:"filename"`
	Duration *int   `json:"duration,omitempty"`
	Size     *int   `json:"size,omitempty"`
}

type ArchiveItemChan struct {
	MessageId int `json:"msg_id,omitempty"`
}

type Item struct {
	Uid      string `json:"uid"`
	Basename string `json:"base_name"`
	Filename string `json:"file_name"`

	Date     string `json:"date"`
	Provider string `json:"provider"`
	ShowName string `json:"show_name"`

	ShowTitle string `json:"show_title"`
	EpTitle   string `json:"ep_title"`

	Artists []string `json:"artists"`

	Size      *int `json:"size"`
	MessageId *int `json:"message_id"`
	Duration  *int `json:"duration"`
}
