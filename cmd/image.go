package cmd

import (
	"github.com/pgeowng/japoto-dl/internal/model"
	"github.com/pgeowng/japoto-dl/internal/provider/hibiki"
	"github.com/pgeowng/japoto-dl/internal/provider/onsen"
	"github.com/pgeowng/japoto-dl/internal/repo/dl"
	"github.com/pgeowng/japoto-dl/internal/repo/status"
	"github.com/pgeowng/japoto-dl/internal/repo/wd"
	"github.com/pgeowng/japoto-dl/internal/usecase"
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
	providers := []usecase.Provider{
		onsen.NewOnsen(),
		hibiki.NewHibiki(),
	}
	status := &status.ErrorPrintLine{}
	wd1 := wd.NewWd("./", "")

	MapShow(d, providers, status, func(show model.Show) error {
		if err := show.LoadImage(d, wd1); err != nil {
			return err
		}
		return nil
	})

}
