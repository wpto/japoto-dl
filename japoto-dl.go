package main

import (
	"github.com/pgeowng/japoto-dl/cmd"
	"github.com/pgeowng/japoto-dl/cmd/html"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "japoto-dl"}
	root.AddCommand(cmd.ListCmd())
	root.AddCommand(cmd.DownloadCmd())
	root.AddCommand(cmd.ImageCmd())
	root.AddCommand(html.Cmd())
	root.Execute()
}
