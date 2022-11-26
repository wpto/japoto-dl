package usecase

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/provider/hibiki"
	"github.com/pgeowng/japoto-dl/internal/provider/onsen"
	"github.com/pgeowng/japoto-dl/internal/repo/dl"
	"github.com/pgeowng/japoto-dl/internal/repo/status"
	"github.com/spf13/cobra"
)

type ListShows struct {
	d          DL
	showMapper ShowMapper
}

func NewListShows() *ListShows {
	d := dl.NewGrequests()
	providers := []Provider{
		onsen.NewOnsen(),
		hibiki.NewHibiki(),
	}

	return &ListShows{
		d: d,
		showMapper: ShowMapper{
			d:         d,
			providers: providers,
			status:    &status.ErrorPrintLine{},
		},
	}
}

func (ls *ListShows) Run(cmd *cobra.Command, args []string) {
	ls.showMapper.MapShow(ls.visitShow)
}

func (ls *ListShows) visitShow(show model.Show) error {
	fmt.Println(show.PPrint().String())

	eps, err := show.GetEpisodes(ls.d)
	if err != nil {
		fmt.Printf("GetEpisodes: error=%v\n", err)
	}
	for _, ep := range eps {
		fmt.Println(ep.PPrint().String())
	}
	return nil
}
