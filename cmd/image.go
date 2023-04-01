package cmd

import (
	"log"

	"github.com/pgeowng/japoto-dl/pkg/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/repo/status"
	"github.com/pgeowng/japoto-dl/workdir/wd"
	"github.com/spf13/cobra"
)

func ImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Loads shows",
		Long:  `Loads shows images`,
		Run:   imageRun,
	}

	cmd.Flags().StringSliceVarP(&FilterProviderList, "provider-only", "p", []string{}, "Shows only selected providers")
	cmd.Flags().StringSliceVarP(&FilterShowIdList, "show-only", "s", []string{}, "Shows only selected shows")
	return cmd
}

func imageRun(cmd *cobra.Command, args []string) {

	d := dl.NewGrequests()
	providers := provider.NewProvidersList()
	status := &status.ErrorPrintLine{}
	wd1 := wd.NewWd("./", "")

	sm := &ShowMapper{dl: d, providers: providers, pl: status}
	shows, err := sm.MapShows()
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, show := range shows {
		if err := show.LoadImage(d, wd1); err != nil {
			log.Fatal(err)
			return
		}
	}
}
