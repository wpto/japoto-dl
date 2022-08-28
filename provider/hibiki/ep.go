package hibiki

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider/common"
	"github.com/pkg/errors"
)

type HibikiEpisode struct {
	Id   int    `json:"id"`
	Name string `json:"name"`

	UpdatedAt string `json:"updated_at"`

	Video           *HibikiEpisodeMedia `json:"video"`
	AdditionalVideo *HibikiEpisodeMedia `json:"additional_video"`
}

type HibikiEpisodeMedia struct {
	Id        int `json:"id"`
	MediaType int `json:"media_type"`
	URL       *string

	IsAdditional bool
	epRef        *HibikiEpisode
	showRef      *HibikiShow
}

func (m *HibikiEpisodeMedia) Artists() []string {
	return []string{}
}

func (m *HibikiEpisodeMedia) CanDownload() bool {
	return true
}

func (m *HibikiEpisodeMedia) Date() model.Date {
	const longForm = "2006/01/02 15:04:05"
	t, err := time.Parse(longForm, m.epRef.UpdatedAt)

	if err != nil {
		panic(errors.Wrap(err, "hibiki.epm.date"))
	}

	year, month, day := t.Date()
	return model.Date{
		Year:  year,
		Month: int(month),
		Day:   day,
	}
}

func (m *HibikiEpisodeMedia) EpId() string {
	return fmt.Sprintf("%s/%s", m.showRef.ShowId(), m.EpIdx())
}

func (m *HibikiEpisodeMedia) EpTitle() string {
	backstage := ""
	if m.IsAdditional {
		backstage += " (楽屋裏)"
	}
	return m.epRef.Name + backstage
}

func (m *HibikiEpisodeMedia) IsVideo() bool {
	return m.MediaType != 1
}

func (m *HibikiEpisodeMedia) Show() model.Show {
	v := reflect.ValueOf(m.showRef).Interface()
	c := v.(model.Show)
	return c
}

func (m *HibikiEpisodeMedia) PPrint() model.PPrintRow {
	return model.PPrintRow{
		IsDir:   false,
		CanLoad: m.CanDownload(),
		IsVid:   m.IsVideo(),
		Date:    m.Date(),
		Ref:     m.EpId(),
		Note:    m.EpTitle(),
		Cast:    []string{},
	}
}

func (m *HibikiEpisodeMedia) EpIdx() string {
	return common.EncodeIdx(m.epRef.Id, m.Id)
}

func (m *HibikiEpisodeMedia) PlaylistURL() *string {
	return m.URL
}
