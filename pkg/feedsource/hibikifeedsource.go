package feedsource

import (
	"encoding/json"
	"net/http"
)

type HikibiFeedEntry struct {
	AccessId    string `json:"access_id"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	NameKana    string `json:"name_kana"`
	DayOfWeek   int    `json:"day_of_week"`
	Description string `json:"description"`
	PCImageURL  string `json:"pc_image_url"`
	PCImageInfo struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"pc_image_info"`
	SPImageURL  string `json:"sp_image_url"`
	SPImageInfo struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"sp_image_info"`
	OnAirInformation    string `json:"onair_information"`
	MessageFromURL      string `json:"message_from_url"`
	Email               string `json:"email"`
	NewProgramFlag      bool   `json:"new_program_flg"`
	PickupProgramFlag   bool   `json:"pickup_program_flg"`
	OriginalProgramFlag bool   `json:"original_program_flg"`
	Copyright           string `json:"copyright"`
	Priority            int    `json:"priority"`
	MetaTitle           string `json:"meta_title"`
	MetaKeyword         string `json:"meta_keyword"`
	MetaDescription     string `json:"meta_description"`
	HashTag             string `json:"hash_tag"`
	ShareText           string `json:"share_text"`
	ShareURL            string `json:"share_url"`
	Cast                string `json:"cast"`
	PublishStartAt      string `json:"publish_start_at"`
	PublishEndAt        string `json:"publish_end_at"`
	UpdatedAt           string `json:"updated_at"`
	LatestEpisodeID     int    `json:"latest_episode_id"`
	LatestEpisodeName   string `json:"latest_episode_name"`
	EpisodeUpdatedAt    string `json:"episode_updated_at"`
	UpdateFlag          bool   `json:"update_flg"`
	Episode             struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		MediaType *struct{} `json:"media_type"`
		Video     struct {
			ID       int         `json:"id"`
			Duration json.Number `json:"duration"`
		} `json:"video"`
		AdditionalVideo *struct{} `json:"additional_video"`
		HTMLDescription string    `json:"html_description"`
		LinkURL         string    `json:"link_url"`
		UpdatedAt       string    `json:"updated_at"`
		EpisodePart     *struct{} `json:"episode_part"`
		Chapters        *struct{} `json:"chapters"`
	} `json:"episode"`
	ChapterFlag         bool `json:"chapter_flg"`
	AdditionalVideoFlag bool `json:"additional_video_flg"`
}

/*
 {
   "access_id": "aokumanowtraining",
   "id": 304,
   "name": "青山なぎさ・熊田茜音のnow training！！",
   "name_kana": "アオヤマナギサ・クマダアカネのnow training！！",
   "day_of_week": 1,
   "description": "『青山なぎさ・熊田茜音のnow training！！』\r\n\r\nこの番組は、Apollo Bay所属の青山なぎさと熊田茜音が勢いをつけるため、\r\n2人の新たな一面や魅力を引き出していく企画挑戦ラジオ番組です！！\r\n\r\n2人へのお便りはもちろん、\r\n質問や2人に話して欲しいことなど幅広く募集！！\r\nみなさまからのお便りお待ちし
ます！！\r\n\r\n■更新情報■\r\n毎週月曜日12時 更新\r\n\r\nまた、ニコニコ響ラジオチャンネルではおまけ付きアーカイブを公開中！\r\n本編では聞けないゆったりとした内容でお届けしています！\r\n是非会員になってチェック！！\r\n\r\n▼響ラジオチャンネル（ニコニコ）\r\nhttps://ch.nicovideo.jp/hibiki",
   "pc_image_url": "https://hibikiradiovms.blob.core.windows.net/image/uploads/program_banner/image1/404/b59c7912-ffda-4b38-81c1-dd419892f663.png",
   "pc_image_info": {
     "width": 1000,
     "height": 400
   },
   "sp_image_url": "https://hibikiradiovms.blob.core.windows.net/image/uploads/program_banner/image2/404/thumb_62f7cd1a-8b8b-42c5-9365-5004cfefd3e3.png",
   "sp_image_info": {
     "width": 640,
     "height": 360
   },
   "onair_information": "【毎週月曜日 12:00 更新】\r\n響ラジオステーションにて\r\n11月1日（月）より\r\n配信スタート！！\r\nニコニコにて\r\n会員限定でおまけ＋アーカイブを配信！！",
   "message_form_url": "https://vcms-api.hibiki-radio.jp/inquiries/new?program_id=304",
   "email": "aokuma2nowtraining@hibiki-radio.jp",
   "new_program_flg": false,
   "pickup_program_flg": true,
   "original_program_flg": true,
   "copyright": "©HiBiKi Radio Station",
   "priority": 240,
   "meta_title": "青山なぎさ・熊田茜音のnow training！！",
   "meta_keyword": "青山なぎさ,熊田茜音,nowtraining,響,響ラジオ,響ラジオステーション",
   "meta_description": "",
   "hash_tag": "",
   "share_text": "響ラジオステーションで「青山なぎさ・熊田茜音のnow training！！」を楽しもう!",
   "share_url": "https://bit.ly/3lDhraI",
   "cast": "青山なぎさ, 熊田茜音",
   "publish_start_at": "2021/10/14 15:00:00",
   "publish_end_at": null,
   "updated_at": "2022/05/02 12:09:26",
   "latest_episode_id": 13449,
   "latest_episode_name": "第53回",
   "episode_updated_at": "2022/10/31 12:00:00",
   "update_flg": false,
   "episode": {
     "id": 13449,
     "name": "第53回",
     "media_type": null,
     "video": {
       "id": 16832,
       "duration": 2778.79
     },
     "additional_video": null,
     "html_description": "",
     "link_url": "",
     "updated_at": "2022/10/31 03:23:19",
     "episode_parts": null,
     "chapters": null
   },
   "chapter_flg": false,
   "additional_video_flg": false
 }

*/

type HibikiFeedSource struct {
	cl http.Client
}

func NewHibikiFeedSource() *HibikiFeedSource {
	return &HibikiFeedSource{}
}

func (f *HibikiFeedSource) GetShowList() (shows []Show, err error) {

	req, err := http.NewRequest("GET", "https://vcms-api.hibiki-radio.jp/api/v1//programs?limit=99&page=1", nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := f.cl.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var result []struct {
		AccessID string `json:"access_id"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return
	}

	shows = make([]Show, len(result))
	for i, v := range result {
		shows[i] = Show{ID: v.AccessID}
	}

	return
}
