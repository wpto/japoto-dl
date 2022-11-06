package expanddb

import (
	"fmt"
	"log"
	"regexp"

	"github.com/pgeowng/japoto-dl/html/config"
	"github.com/pgeowng/japoto-dl/html/entity"
	"github.com/pgeowng/japoto-dl/html/store"
	"github.com/pgeowng/japoto-dl/html/types"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expanddb",
		Short: "Parses db and add new info",
		Long:  `expands channel info by adding new parsed values needed later`,
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {

	store := store.NewFileStore(config.FileStorePath)

	entries := store.Read()
	entries = ExtendContent(entries)
	// entries = ExtendPerformers(entries)
	store.Write(entries)
	log.Println("Done")
}

func ExtendContent(eps []types.Entry) []types.Entry {
	for idx := range eps {
		//eps[idx].Provider = "unknown"
		//eps[idx].Date = "000000"
		//eps[idx].ShowId = "unknown"

		fmt.Printf("%d: %+v\n", idx, eps[idx])
		info, err := GuessMeta(eps[idx].Filename)
		if err != nil {
			log.Fatal(err)
		}

		eps[idx].Date = info.Date
		eps[idx].ShowId = info.ShowId
		eps[idx].Provider = info.Provider

		if len(eps[idx].Title) == 0 {
			eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].ShowId)
		}

		prefixTitleRE := regexp.MustCompile(`^(\d{6})`)
		match := prefixTitleRE.FindStringSubmatch(eps[idx].Title)
		if len(match) == 0 {
			eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].Title)
		}

		eps[idx].URL = config.ChannelPrefix + fmt.Sprint(eps[idx].MessageId)
		eps[idx].DurationHuman = entity.FormatDurationHuman(eps[idx].Duration)
		eps[idx].SizeHuman = entity.FormatSizeHuman(eps[idx].Size)
	}

	return eps
}

func ExtendPerformers(eps []types.Entry) []types.Entry {
	for idx := range eps {
		info, err := GuessPerformers(eps[idx].Performer)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("ok: %v\n", info)
		}
	}
	return eps
}
