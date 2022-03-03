package cmd

import "github.com/spf13/cobra"

func ImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Loads shows",
		Long:  `Loads shows images`,
		Run:   imageRun,
	}

	return cmd
}

func imageRun(cmd *cobra.Command, args []string) {

}
