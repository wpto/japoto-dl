package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crawshaw.io/sqlite/sqlitex"
	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/pkg/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/repo/archive"
	status1 "github.com/pgeowng/japoto-dl/repo/status"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/wd"
)

const (
	feedCheckTime = 5 * time.Second
	showCheckTime = 5 * time.Second
)

func main() {
	ctx := context.Background()

	app, err := NewApp()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
		return
	}

}

type App struct {
	db     *sqlitex.Pool
	arch   *archive.ArchiveRepo
	loader model.Loader
	hibiki *provider.Hibiki
	onsen  *provider.Onsen

	muxInput  chan MuxJob
	muxWorker *MuxWorker
}

func NewApp() (*App, error) {
	db, err := archive.CreateDB("./jwd.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	arch, err := archive.NewRepo(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	loader := dl.NewGrequests()
	hibiki := provider.NewHibiki(loader)
	onsen := provider.NewOnsen(loader)

	status := status1.New(os.Stdout)
	muxInput := make(chan MuxJob, 0)
	muxWorker := NewMuxWorker(status, muxInput)

	return &App{
		db:        db,
		arch:      arch,
		loader:    loader,
		hibiki:    hibiki,
		onsen:     onsen,
		muxInput:  muxInput,
		muxWorker: muxWorker,
	}, nil
}
func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	onsenTicker := time.NewTicker(feedCheckTime)
	hibikiTicker := time.NewTicker(feedCheckTime)
	// onsenShowTicker := time.NewTicker(showCheckTime)
	hibikiShowTicker := time.NewTicker(showCheckTime)
	downloadTicker := time.NewTicker(showCheckTime)

	_ = onsenTicker
	_ = hibikiTicker
	_ = hibikiShowTicker
	_ = downloadTicker

	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-stopChan
		cancel()
	}()

	if err := a.arch.Exec(ctx, createShowsTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createPersonaTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createPersonaRoleTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createPersonaShowTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createShowKeyTable); err != nil {
		fmt.Println(err)
		return err
	}

	q := `insert into show_keys(id, label)
	      values (1, 'title'),
	             (2, 'poster')
	      on conflict do nothing`
	if err := a.arch.Exec(ctx, q); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createShowKeyValueTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createURLBank); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createEpisodeTable); err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createDownloadJobTable); err != nil {
		fmt.Println(err)
		return err
	}

	_ = downloadTicker
	fmt.Println("app is running...")
	fmt.Println("run download job")
	// if err := a.RunDownloadWorker(ctx); err != nil {
	// 	fmt.Println(err)
	// }
	if err := a.RunOnsenGetFeed(ctx); err != nil {
		fmt.Println(err)
	}
	for {
		// if err := a.RunGetRecentShows(ctx, a.onsen); err != nil {
		// 	fmt.Println(err)
		// }
		downloadWorker := NewDownloadWorker(a.arch).WithPostDownloadHook(func(wdhls *workdir.Workdir, wd *wd.Wd) {
			a.muxInput <- MuxJob{wdhls: wdhls, wd: wd}
		})
		if err := downloadWorker.RunDownloadWorker(ctx); err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Second)
		select {
		// case <-onsenTicker.C:
		// 	if err := a.RunOnsenGetFeed(ctx); err != nil {
		// 		fmt.Println(err)
		// 	}
		// case <-hibikiTicker.C:
		// 	if err := a.RunHibikiGetFeed(ctx); err != nil {
		// 		fmt.Println(err)
		// 	}
		// case <-onsenShowTicker.C:
		// 	if err := a.RunGetRecentShows(ctx, a.onsen); err != nil {
		// 		fmt.Println(err)
		// 	}
		// case <-hibikiShowTicker.C:
		// 	if err := a.RunGetRecentShows(ctx, a.hibiki); err != nil {
		// 		fmt.Println(err)
		// 	}
		case <-stopChan:
			fmt.Println("closing app by interrupt")
			cancel()
			return fmt.Errorf("interrupt")
		default:
		}
	}
}

const createShowsTable = `create table if not exists shows (
	id integer primary key autoincrement,
	source text not null,
	show_name text not null,
	inserted_at text default current_timestamp not null,
	worker_id string default '' not null,
	state string default '' not null,
	updated_at text default current_timestamp not null,
	unique (source, show_name)
)`

const createPersonaTable = `create table if not exists persona (
  id integer primary key autoincrement,
  name text not null,
  inserted_at text default current_timestamp not null,
  unique (name)
)`

const createPersonaRoleTable = `create table if not exists persona_role (
  id integer primary key autoincrement,
  name text not null,
  show_id integer not null,
  persona_id integer not null,
  inserted_at text default current_timestamp not null,

  foreign key (show_id) references shows (id),
  foreign key (persona_id) references persona (id),
  unique (show_id, persona_id, name)
)`

// TODO: support persona that was in show but currency not presented
//
//	probably episode based
const createPersonaShowTable = `create table if not exists persona_show (
  id integer primary key autoincrement,
  show_id integer not null,
  persona_id integer not null,
  inserted_at text default current_timestamp not null,

  foreign key (show_id) references shows (id),
  foreign key (persona_id) references persona (id),
  unique (show_id, persona_id)
)`

const createShowKeyTable = `create table if not exists show_keys (
	id integer primary key autoincrement,
	label text not null,
	inserted_at text default current_timestamp not null,
	unique (label)
)`

const createShowKeyValueTable = `create table if not exists show_values (
  id integer primary key autoincrement,
  show_id integer not null,
  key integer not null,
  value text not null,
  inserted_at text default current_timestamp not null,

  foreign key (show_id) references shows (id),
  foreign key (key) references show_keys (id)
)`

const createURLBank = `create table if not exists url_bank (
	id integer primary key autoincrement,
	url text not null,
	inserted_at text default current_timestamp not null,
	unique (url)
)`

const createEpisodeTable = `create table if not exists show_episodes (
	id integer primary key autoincrement,
	show_id integer not null,
  title text not null,
  date_day integer not null,
  date_month integer not null,
  date_year integer not null,
  inserted_at text default current_timestamp not null,

  foreign key (show_id) references shows (id),
  unique (show_id, title, date_day, date_month, date_year)
)`

const createPersonaGuestTable = `create table if not exists persona_guest(
  id integer primary key autoincrement

  ep_id integer not null,
  persona_id integer not null,
  inserted_at text default current_timestamp not null,

  foreign key (ep_id) references show_episodes (id),
  foreign key (persona_id) references persona (id),
  unique (show_id, persona_id)
)`

const createDownloadJobTable = `create table if not exists download_jobs(
  id integer primary key autoincrement,

  ep_id integer not null,
  playlist_url integer not null,
  image_url integer not null,

  inserted_at text default current_timestamp not null,

  worker_id text default '' not null,
  updated_at text default current_timestamp not null,
  status text default '' not null,

  foreign key (ep_id) references show_episodes(id),
  foreign key (playlist_url) references url_bank(id),
  foreign key (image_url) references url_bank(id),

  unique(ep_id)
)`

func (a *App) RunHibikiGetFeed(ctx context.Context) error {
	fmt.Println("requesting hibiki feed")
	access, err := a.hibiki.GetFeedW(a.loader)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createShowsTable); err != nil {
		fmt.Println(err)
		return err
	}

	for _, it := range access {
		err := a.arch.InsertShow(ctx, "hibiki", it.AccessId)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println(access)

	return nil
}

func (a *App) RunOnsenGetFeed(ctx context.Context) error {
	fmt.Println("requesting onsen feed")
	access, err := a.onsen.GetFeedW()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := a.arch.Exec(ctx, createShowsTable); err != nil {
		fmt.Println(err)
		return err
	}

	for _, it := range access {
		err := a.arch.InsertShow(ctx, "onsen", it.DirectoryName)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println(access)

	return nil
}
