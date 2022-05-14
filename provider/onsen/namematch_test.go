package onsen

import "testing"

type Case struct {
	url    string
	showId string
	err    bool
	yy     int
	mm     int
	dd     int
	num    string
}

func TestExtract(t *testing.T) {
	cases := []Case{
		{"https://onsen-ma3phlsvod.sslcs.cdngc.net/onsen-ma3pvod/_definst_/202003/ore.ski200326l3md-16.mp4/playlist.m3u8",
			"ore.ski",
			false,
			2020, 3, 26,
			"16",
		},
		{"/202201/fujita-pLs1S-04.mp4", "fujita-p",
			false,
			2022, 1, -1,
			"04",
		},
		{"/202111/gurepap21114av8u-55.mp4", "gurepap",
			false,
			2021, 11, 4,
			"55",
		},
		{"/202202/maho7220209sbz4-13.mp4", "maho7",
			false,
			2022, 2, 9,
			"13",
		},
		{"/202112/g123211228iwb7-06.mp4", "g123",
			false,
			2021, 12, 28,
			"06",
		},
		{"/202201/86220121z9nm-23.mp4", "86",
			false,
			2022, 1, 21,
			"23",
		},
		{"/202202/fujita-p2200204nNE8-91.mp4", "fujita-p",
			false,
			2022, 2, 4,
			"91",
		},
		{"/202110/fujita211022adv7yd4g-97.mp4", "fujita",
			false,
			2021, 10, 22,
			"97",
		},
		{"/202202/aniradiaward220204kd3y-sp.mp4", "aniradiaward",
			false,
			2022, 2, 4,
			"sp",
		},
		{"/202109/fuchigamimai210921gvx0-sp2.mp4", "fuchigamimai",
			false,
			2021, 9, 21,
			"sp2",
		},
		{"/202202/techno-roid220220fryj-10.mp4", "techno-roid",
			false,
			2022, 2, 20,
			"10",
		},
		{"/202202/gaikotukishi220218ue7g-02.mp4/playlist.m3u8", "gaikotsukishi",
			false,
			2022, 2, 18,
			"02",
		},
		{"/202202/d220226l3md-16-sp3.mp4/playlist.m3u8", "d",
			false,
			2022, 2, 26,
			"16-sp3",
		},
		{"/202202/nkm220208dx7e-41-2.mp4/playlist.m3u8", "nkm",
			false,
			2022, 2, 8,
			"41-2",
		},
		{"https://onsen-ma3phlsvod.sslcs.cdngc.net/onsen-ma3pvod/_definst_/archive/survey_01.mp4/playlist.m3u8", "survey",
			false,
			-1, -1, -1,
			"",
		},
	}

	for _, c := range cases {
		result, err := Extract(c.url, c.showId)
		if c.err && err == nil {
			t.Logf("(%s, %s) expected err, got %#v", c.url, c.showId, result)
			t.Fail()
		} else if !c.err && err != nil {
			t.Logf("(%s, %s) expected result, got err %#v", c.url, c.showId, err)
			t.Fail()
		} else {
			flag := false
			if result.DateY != c.yy {
				t.Logf("wrong year %d != %d", result.DateY, c.yy)
				flag = true
			}
			if result.DateM != c.mm {
				t.Logf("wrong month %d != %d", result.DateM, c.mm)
				flag = true
			}
			if result.Num != c.num {
				t.Logf("wrong num %s != %s", result.Num, c.num)
				flag = true
			}
			if result.DateD != c.dd {
				t.Logf("wrong day %v != %v", result.DateD, c.dd)
				flag = true
			}
			if flag {
				t.Logf("for call (%s, %s)", c.url, c.showId)
				t.Fail()
			}
		}
	}
}
