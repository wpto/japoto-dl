package archive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"github.com/pgeowng/japoto-dl/internal/model"
)

var (
	ErrDuplicate = errors.New("record already exists")
)

type Archive interface {
	Migrate(archive *sqlitex.Pool) error
	IsLoaded(archive *sqlitex.Pool, key string) (ok bool, err error)
	SetLoaded(archive *sqlitex.Pool, key string, status bool) error
	Create(archive *sqlitex.Pool, key string, status bool, archiveItem model.ArchiveItem) error
}

type ArchiveRepo struct{}

func CreateDB(filename string) (pool *sqlitex.Pool, err error) {
	pool, err = sqlitex.Open(filename, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE, 10)
	if err != nil {
		err = fmt.Errorf("CreateDB: sqlitex.Open: %w", err)
		return
	}

	return
}

func NewRepo() (archive *ArchiveRepo, err error) {
	return &ArchiveRepo{}, nil
}

func (r *ArchiveRepo) Migrate(pool *sqlitex.Pool) (err error) {
	conn := pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	query := `CREATE TABLE IF NOT EXISTS history(key TEXT PRIMARY KEY, status INTEGER, data json);`
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.Migrate: Finalize: %w", err)
				return
			}
		}
	}()
	if _, err = stmt.Step(); err != nil {
		return
	}

	query = `create table if not exists episodes(
		key string primary key,
		date string,
		source string,
		show_id string,
		ep_id string,
		ep_title string,
		artists string
	);`

	stmt = conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.Migrate: Finalize: %w", err)
				return
			}
		}
	}()
	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

type HistoryItem struct {
	Url string `json:"string"`
}

type ArchiveEntryStatus int

const (
	Unknown ArchiveEntryStatus = iota
	NotExists
	NotLoaded
	Loaded
)

func (r *ArchiveRepo) IsLoaded(pool *sqlitex.Pool, key string) (status ArchiveEntryStatus, err error) {
	conn := pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	query := `SELECT status FROM history WHERE key = $key;`

	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.IsLoaded: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	rowReturned, err := stmt.Step()
	if err != nil {
		err = fmt.Errorf("ArchiveRepo.IsLoaded: Step: %w", err)
		return
	}

	if rowReturned {
		status = ArchiveEntryStatus(stmt.GetInt64("status"))
	} else {
		status = NotExists
	}

	return
}

func (r *ArchiveRepo) SetStatus(pool *sqlitex.Pool, key string, status ArchiveEntryStatus) (err error) {
	conn := pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	query := `INSERT INTO history(key, status, data) VALUES ($key, $status, '{}') ON CONFLICT (key) DO UPDATE SET status = $status;`
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.SetStatus: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	stmt.SetInt64("$status", int64(status))
	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

func (r *ArchiveRepo) Create(pool *sqlitex.Pool, key string, status ArchiveEntryStatus, archiveItem model.ArchiveItem) (err error) {
	bytes, err := json.Marshal(archiveItem)
	if err != nil {
		return fmt.Errorf("ArchiveRepo.Create: can't marshal archiveItem: %w", err)
	}

	conn := pool.Get(context.TODO())
	if conn == nil {
		return fmt.Errorf("ArchiveRepo.Create: couldn't get connection for pool")
	}
	defer pool.Put(conn)

	query := "INSERT INTO history(key, status, data) values($key, $status, $data);"
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.Create: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	stmt.SetInt64("$status", int64(status))
	stmt.SetText("$data", string(bytes))

	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

// func (r *ArchiveRepo) createNewItem(pool *sqlitex.Pool, key string, status ArchiveEntryStatus, archiveItem model.ArchiveItem) (err error) {
// 	conn := pool.Get(context.TODO())
// 	if conn == nil {
// 		return fmt.Errorf("ArchiveRepo.Create: couldn't get connection for pool")
// 	}
// 	defer pool.Put(conn)

// 	query := "insert into history_v2(archive_key, date, source, show_id, show_title, ep_id, ep_title, artists, filename, duration, size, msg_id) values ($archive_key, $date, $source, $show_id, $show_title, $ep_id, $ep_title, $artists, $filename, $duration, $size, $msg_id);"

// 	stmt := conn.Prep(query)
// 	defer func() {
// 		if err == nil {
// 			if err = stmt.Finalize(); err != nil {
// 				err = fmt.Errorf("ArchiveRepo.Create: Finalize: %w", err)
// 				return
// 			}
// 		}
// 	}()

// 	stmt.SetText("$key", key)
// 	stmt.SetInt64("$status", int64(status))
// 	stmt.SetText("$data", string(bytes))
// 	if _, err = stmt.Step(); err != nil {
// 		return
// 	}

// 	return
// }

// func (a *ArchiveRepo) CreateEpisode(pool *sqlitex.Pool, key string, ep model.ArchiveEpisode) (err error) {
// 	conn := pool.Get(context.TODO())
// 	if conn == nil {
// 		return fmt.Errorf("ArchiveRepo.CreateEpisode: couldn't get connection for pool")
// 	}
// 	defer pool.Put(conn)

// 	query := "insert into episodes(key, date, source, show_id, ep_id, ep_title, artists) values "
// 	query += "($key, $date, $source, $show_id, $ep_id, $ep_title, $artists);"

// 	stmt := conn.Prep(query)
// 	defer func() {
// 		if err == nil {
// 			if err = stmt.Finalize(); err != nil {
// 				err = fmt.Errorf("ArchiveRepo.Create: Finalize: %w", err)
// 				return
// 			}
// 		}
// 	}()

// 	stmt.SetText("$key", key)
// 	stmt.SetText("$date", ep.Date)
// 	stmt.SetText("$source", ep.Source)
// 	stmt.SetText("$show_id", ep.ShowID)
// 	stmt.SetText("$ep_id", ep.EpID)
// 	stmt.SetText("$ep_title", ep.EpTitle)
// 	stmt.SetText("$artists", strings.Join(ep.Artists, ";"))

// 	if _, err = stmt.Step(); err != nil {
// 		return
// 	}

// 	return
// }
