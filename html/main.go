package main

import (
	"github.com/pgeowng/japoto/expanddb"
	"github.com/pgeowng/japoto/pagegen"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(pagegen.Cmd())
	rootCmd.AddCommand(expanddb.Cmd())
	rootCmd.Execute()
}
