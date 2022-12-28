package feedsource

type Show struct {
	ID string
}

type FeedSource interface {
	GetShowList() ([]Show, error)
}
