package cmd

import (
	"fmt"

	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/model"
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

	MapShow(d, providers, status, func(show model.Show) error {
		fmt.Println(show.PPrint().String())
		eps, err := show.GetEpisodes(d)
		if err != nil {
			fmt.Printf("GetEpisodes: error=%v\n", err)
		}
		for _, ep := range eps {
			fmt.Println(ep.PPrint().String())
		}
		return nil
	})
}
