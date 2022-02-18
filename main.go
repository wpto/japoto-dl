package main

import (
	"fmt"
	"log"

	"github.com/levigross/grequests"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pkg/errors"
)

type Loader struct{}

// options should me separated from code...
func (loader *Loader) Text(url string, opts *grequests.RequestOptions) (*string, error) {
	res, err := grequests.Get(url, opts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}
	str := res.String()
	return &str, nil
}

func (loader *Loader) JSON(url string, dest interface{}, opts *grequests.RequestOptions) error {
	res, err := grequests.Get("https://onsen.ag/web_api/programs/", opts)
	if err != nil {
		return errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	// mapObj := make([]map[string]interface{}, 0)
	err = res.JSON(dest)
	if err != nil {
		return errors.Wrap(err, "get.json")
	}

	return nil
}

func (loader *Loader) Raw(url string, opts *grequests.RequestOptions) ([]byte, error) {
	res, err := grequests.Get(url, opts)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	fmt.Printf("%s: %v\n", url, res.StatusCode)
	if !res.Ok {
		return nil, errors.Errorf("get: bad status code(%d): %s", res.StatusCode, url)
	}

	return res.Bytes(), nil
}

func main() {
	loader := &Loader{}
	wd := workdir.NewWorkdir("./.cache/", "fasqwge")
	tasks := tasks.NewTasks(wd)
	providersInfo := provider.NewProvidersInfo(loader)
	providers := provider.NewProviders(loader, tasks)

	shows, err := providersInfo.Onsen.GetFeed()
	if err != nil {
		log.Fatal(err)
		return
	}

	// fmt.Println("len", len(shows))
	for _, show := range shows {
		// fmt.Println(show.ShowId())
		eps := show.GetEpisodes()
		// fmt.Println("ep.len", len(eps))
		for _, ep := range eps {
			url := ep.PlaylistUrl()
			if url != nil {
				err := providers.Onsen.Download(*url)
				fmt.Println(err)
				return
			}
			fmt.Printf("%s ", ep.EpTitle())
		}
	}

	// log.Printf("%#v\n", result)

}
