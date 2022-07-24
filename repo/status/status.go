package status

import "github.com/pgeowng/japoto-dl/internal/types"

type LoadStatusImpl struct {
	Provider  string
	EpisodeID string
	total     int
	loaded    int
}

func NewLoadStatus(provider string, episodeID string) types.LoadStatus {
	return &LoadStatusImpl{Provider: provider, EpisodeID: episodeID}
}

func (st *LoadStatusImpl) Inc(step int) {
	st.loaded += step
	// TODO: print
}

func (st *LoadStatusImpl) Total(total int) {
	// TODO: total already set
	st.total = total
	// TODO: print
}
