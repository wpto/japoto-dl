package cmd

import (
	"fmt"
	"sort"

	"github.com/pgeowng/japoto-dl/client/hibiki"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func ShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows",
		Short: "Shows",
		Long:  `Shows`,
		Run:   run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	c := hibiki.New()
	showList, err := c.GetShowList()
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(showList, func(i, j int) bool {
		return showList[i] < showList[j]
	})

	cellWidth := 0
	for i := 0; i < len(showList); i++ {
		if len(showList[i]) > cellWidth {
			cellWidth = len(showList[i])
		}
	}

	width, _, err := terminal.GetSize(0)
	if err != nil {
		width = 0
	}

	cellWidth += 2

	columns := width / cellWidth
	if columns > 5 {
		columns = 5
	}

	// cell rounded up to nearest multiple of columns
	columnSize := (len(showList) + columns - 1) / columns

	for i := 0; i < columnSize; i++ {
		for j := 0; j < columns; j++ {
			idx := i + j*columnSize
			if idx < len(showList) {
				fmt.Printf("%-*s", cellWidth, showList[idx])
			}
		}
		fmt.Println()
	}
}
