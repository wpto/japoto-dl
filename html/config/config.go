package config

import (
	"fmt"
	"time"
)

var Dest = "./public"
var Static = "./html/static"
var TemplateDir = "./html/template/"
var RecentLimit = 30
var FileStorePath = "../tgchan/"
var ChannelPrefix = "https://t.me/japoto/"
var PublicURL = "https://pgeowng.github.io/japoto"

var day = []string{"日", "月", "火", "水", "木", "金", "土"}
var jst = time.FixedZone("UTC+9", 9*60*60)
var now = time.Now().In(jst)
var CreateTime = fmt.Sprintf(now.Format("060102 %s 15:04 JST"), day[now.Weekday()])

func init() {
	fmt.Println(CreateTime)
}
