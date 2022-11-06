package printers

import (
	"sort"

	"github.com/pgeowng/japoto-dl/html/types"
)

func FilterEntries(entries []types.Entry, cond func(entry types.Entry) bool) []types.Entry {
	filtered := make([]types.Entry, 0)
	for _, ep := range entries {
		if cond(ep) {
			filtered = append(filtered, ep)
		}
	}
	return filtered
}

func UniqueRecentShows(entries []types.Entry) []types.Entry {

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].MessageId > entries[j].MessageId
	})

	filtered := make([]types.Entry, 0)
	nameSet := make(map[string]bool)
	for _, ep := range entries {
		name := ep.ShowId
		if _, ok := nameSet[name]; !ok {
			filtered = append(filtered, ep)
			nameSet[name] = true
		}
	}

	return filtered
}
