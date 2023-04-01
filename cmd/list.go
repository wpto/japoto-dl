package cmd

import (
	"fmt"
	"log"

	"github.com/pgeowng/japoto-dl/pkg/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/repo/status"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List shows",
		Long:  `List current shows content`,
		Run:   listRun,
	}

	cmd.Flags().StringSliceVarP(&FilterProviderList, "provider-only", "p", []string{}, "Shows only selected providers")
	cmd.Flags().StringSliceVarP(&FilterShowIdList, "show-only", "s", []string{}, "Shows only selected shows")
	return cmd
}

func listRun(cmd *cobra.Command, args []string) {
	d := dl.NewGrequests()
	providers := provider.NewProvidersList()
	status := &status.ErrorPrintLine{}

	sm := &ShowMapper{dl: d, providers: providers, pl: status}
	shows, err := sm.MapShows()
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, show := range shows {
		fmt.Println(show.PPrint().String())
		eps, err := show.GetEpisodes(d)
		if err != nil {
			fmt.Printf("GetEpisodes: error=%v\n", err)
		}

		for _, ep := range eps {
			fmt.Println(ep.PPrint().String())
		}
	}
}
