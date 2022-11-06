package archive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pgeowng/japoto-dl/model"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
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

var ErrNoItem = errors.New("no item found")

func (r *ArchiveRepo) GetArchiveEntry(pool *sqlitex.Pool, key string) (item model.ArchiveItem, err error) {

	conn := pool.Get(context.TODO())
	if err != nil {
		err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: couldn't get connection for pool: %w", err)
		return
	}
	defer pool.Put(conn)

	query := `SELECT data FROM history WHERE key = $key;`
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)

	hasRow, err := stmt.Step()
	if err != nil {
		err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: Step: %w", err)
		return
	}

	if !hasRow {
		err = ErrNoItem
		return
	}

	itemBody := stmt.GetText("data")

	err = json.Unmarshal([]byte(itemBody), &item)
	if err != nil {
		err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: can't unmarshal: %w", err)
		return
	}

	return
}

func (r *ArchiveRepo) GetArchiveEntries(pool *sqlitex.Pool) (items []model.ArchiveItem, err error) {

	conn := pool.Get(context.TODO())
	if err != nil {
		err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: couldn't get connection for pool: %w", err)
		return
	}
	defer pool.Put(conn)

	query := `SELECT data FROM history order by key`
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err = stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: Finalize: %w", err)
				return
			}
		}
	}()

	for {

		var hasRow bool

		hasRow, err = stmt.Step()
		if err != nil {
			err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: Step: %w", err)
			return
		}

		if !hasRow {
			if len(items) == 0 {
				err = ErrNoItem
			}
			return
		}

		var item model.ArchiveItem
		itemBody := stmt.GetText("data")

		err = json.Unmarshal([]byte(itemBody), &item)
		if err != nil {
			err = fmt.Errorf("ArchiveRepo.GetArchiveEntry: can't unmarshal: %w", err)
			return
		}

		items = append(items, item)
	}
}
