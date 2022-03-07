package cmd

import (
	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/provider"
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
	pl := &ErrorPrintLine{}
	wd1 := wd.NewWd("./", "")

	MapShow(d, providers, pl, func(show model.Show) error {
		if err := show.LoadImage(d, wd1); err != nil {
			return err
		}
		return nil
	})

}
