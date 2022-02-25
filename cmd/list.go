package cmd

import (
	"fmt"
	"log"

	"github.com/adios/pprint"
	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/spf13/cobra"
)

var ShowName string = ""

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List shows",
		Long:  `List current shows content`,
		Run:   listRun,
	}

	cmd.Flags().StringVarP(&ShowName, "show-only", "s", "", "Show only <name> show")

	return cmd
}

type Attrs struct {
	IsDir   bool
	CanLoad bool
	IsVid   bool
}

func (a Attrs) String() string {
	result := ""

	if a.IsDir {
		result += "d"
	} else {
		result += "-"
	}
	if a.CanLoad {
		result += "l"
	} else {
		result += "-"
	}
	if a.IsVid {
		result += "v"
	} else {
		result += "-"
	}

	return result
}

func NewAttrs(dir, load, vid bool) Attrs {
	return Attrs{dir, load, vid}
}

func listRun(cmd *cobra.Command, args []string) {
	pp := pprint.NewNode(
		pprint.WithColumns(
			pprint.NewColumn(),
			pprint.NewColumn(),
			pprint.NewColumn(pprint.WithLeftAlignment()),
			pprint.NewColumn(pprint.WithLeftAlignment()),
		),
	)

	d := dl.NewGrequests()
	prov := provider.NewProviders()
	shows, err := prov.Onsen.GetFeed(d)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, show := range shows {
		if ShowName != "" && show.ShowId() != ShowName {
			continue
		}
		eps := show.GetEpisodes()
		_, err := pp.Push(NewAttrs(true, false, false).String(), "", show.ShowId(), show.ShowTitle())
		if err != nil {
			fmt.Println(err)
		}
		for _, ep := range eps {
			url := ep.PlaylistUrl()
			date, err := ep.Date()
			datestr := "------"
			if err != nil {
				fmt.Println(err)
			} else {
				datestr = date.String()
			}

			_, err = pp.Push(NewAttrs(false, url != nil, false), datestr, ep.EpId(), ep.EpTitle())
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	pprint.Print(pp)

}
