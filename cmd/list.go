package cmd

import (
	"github.com/pgeowng/japoto-dl/internal/usecase"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	uc := usecase.NewListShows()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List shows",
		Long:  `List current shows content`,
		Run:   uc.Run,
	}

	cmd.Flags().StringSliceVarP(&FilterProviderList, "provider-only", "p", []string{}, "Shows only selected providers")
	cmd.Flags().StringSliceVarP(&FilterShowIdList, "show-only", "s", []string{}, "Shows only selected shows")
	return cmd
}
