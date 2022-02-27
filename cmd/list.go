package cmd

import (
	"fmt"
	"log"

	"github.com/adios/pprint"
	"github.com/pgeowng/japoto-dl/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/spf13/cobra"
)

var ShowName string = ""

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List shows",
		Long:  `List current shows content`,
		Run:   listRun,
	}

	cmd.Flags().StringVarP(&ShowName, "show-only", "s", "", "Show only <name> show")

	return cmd
}

func listRun(cmd *cobra.Command, args []string) {

	d := dl.NewGrequests()
	prov := provider.NewProviders()
	shows, err := prov.Hibiki.GetFeed(d)
	if err != nil {
		log.Fatal(err)
		return
	}

	// result := []Row{}

	for _, showAccess := range shows {
		if ShowName != "" && showAccess.ShowId() != ShowName {
			continue
		}

		pp := pprint.NewNode(
			pprint.WithColumns(
				pprint.NewColumn(),
				pprint.NewColumn(),
				pprint.NewColumn(pprint.WithLeftAlignment()),
				pprint.NewColumn(pprint.WithLeftAlignment()),
				pprint.NewColumn(pprint.WithLeftAlignment()),
			),
		)

		show, err := showAccess.GetShow(d)
		if err != nil {
			fmt.Println(err)
			continue
		}

		p := show.PPrint().Pprint()
		_, err = pp.Push(p...)
		if err != nil {
			fmt.Println(err)
		}

		for _, ep := range show.GetEpisodes() {
			p = ep.PPrint().Pprint()
			_, err = pp.Push(p...)
			if err != nil {
				fmt.Println(err)
			}
		}
		pprint.Print(pp)
	}
}
