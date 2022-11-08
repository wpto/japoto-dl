package hibiki

import "testing"

func TestClientGetFeed(t *testing.T) {
	c := New()
	feed, err := c.GetFeed()
	if err != nil {
		t.Error(err)
	}

	if len(feed) == 0 {
		t.Error("feed is empty")
	}
}

func TestClientGetShow(t *testing.T) {
	c := New()
	show, err := c.GetShow("imas_cg")
	if err != nil {
		t.Error(err)
	}

	if show.Episode.Video == nil {
		t.Error("video is nil")
	}
}

func TestClientGetShowList(t *testing.T) {
	c := New()
	showList, err := c.GetShowList()
	if err != nil {
		t.Error(err)
	}

	if len(showList) == 0 {
		t.Error("show list is empty")
	}

	t.Log(showList)
	t.Fail()
}
