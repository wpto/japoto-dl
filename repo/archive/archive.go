package archive

import (
	"encoding/json"
	"errors"
	"fmt"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

const dbFilename = "test.pool"

// func main() {
// 	if err := run(); err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("exit")
// }

// func run() (err error) {

// 	repo := NewArchive(pool)
// 	err = repo.Migrate()
// 	if err != nil {
// 		return
// 	}

// 	_, err = repo.Create(Website{
// 		ID:   0,
// 		Name: "World",
// 		URL:  "http://www.sqlite.org",
// 		Rank: 2,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	_, err = repo.Create(Website{
// 		ID:   0,
// 		Name: "Another",
// 		URL:  "http://www.google.com",
// 		Rank: 1,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	all, err := repo.All()
// 	if err != nil {
// 		return
// 	}

// 	fmt.Println(all)

// 	return
// }

var (
	ErrDuplicate = errors.New("record already exists")
)

type Archive interface {
	Migrate(archive *sqlitex.Pool) error
	IsLoaded(archive *sqlitex.Pool, key string) (ok bool, err error)
	SetLoaded(archive *sqlitex.Pool, key string, loaded bool) error
	Create(archive *sqlitex.Pool, key string, loaded bool, archiveItem ArchiveItem) error
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
	conn := pool.Get(nil)
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	// DROP TABLE IF EXISTS entries;
	// CREATE TABLE IF NOT EXISTS entries(
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT NOT NULL UNIQUE,
	// 	key TEXT NOT NULL,
	// 	rank INTEGER NOT NULL
	// );
	query := `DROP TABLE IF EXISTS history;`

	stmt := conn.Prep(query)
	if _, err = stmt.Step(); err != nil {
		return
	}

	query = `CREATE TABLE IF NOT EXISTS history(key TEXT PRIMARY KEY, loaded INTEGER, data json);`
	stmt = conn.Prep(query)
	defer func() {
		if err == nil {
			if err := stmt.Finalize(); err != nil {
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

func (r *ArchiveRepo) IsLoaded(pool *sqlitex.Pool, key string) (ok bool, err error) {
	conn := pool.Get(nil)
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	query := `SELECT EXISTS(SELECT 1 FROM history WHERE key = $key AND loaded = TRUE);`

	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err := stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.IsLoaded: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	if _, err := stmt.Step(); err != nil {
		return false, err
	}

	columnName := stmt.ColumnName(0)
	result := stmt.GetInt64(columnName)

	return result > 0, nil
}

func (r *ArchiveRepo) SetLoaded(pool *sqlitex.Pool, key string, loaded bool) (err error) {
	conn := pool.Get(nil)
	if conn == nil {
		return
	}
	defer pool.Put(conn)

	query := `INSERT INTO history(key, loaded, data) VALUES ($key, $loaded, '{}') ON CONFLICT (key) DO UPDATE SET loaded = $loaded;`
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err := stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.SetLoaded: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	stmt.SetBool("$loaded", loaded)
	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

func (r *ArchiveRepo) Create(pool *sqlitex.Pool, key string, loaded bool, archiveItem ArchiveItem) (err error) {
	bytes, err := json.Marshal(archiveItem)
	if err != nil {
		return fmt.Errorf("ArchiveRepo.Create: can't marshal archiveItem: %w", err)
	}

	conn := pool.Get(nil)
	if conn == nil {
		return fmt.Errorf("ArchiveRepo.Create: couldn't get connection for pool")
	}
	defer pool.Put(conn)

	query := "INSERT INTO history(key, loaded, data) values($key, $loaded, $data);"
	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if err := stmt.Finalize(); err != nil {
				err = fmt.Errorf("ArchiveRepo.Create: Finalize: %w", err)
				return
			}
		}
	}()

	stmt.SetText("$key", key)
	stmt.SetBool("$loaded", loaded)
	stmt.SetText("$data", string(bytes))

	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

type Item struct {
	Uid      string `json:"uid"`
	Basename string `json:"base_name"`
	Filename string `json:"file_name"`

	Date     string `json:"date"`
	Provider string `json:"provider"`
	ShowName string `json:"show_name"`

	ShowTitle string `json:"show_title"`
	EpTitle   string `json:"ep_title"`

	Artists []string `json:"artists"`

	Size      *int `json:"size"`
	MessageId *int `json:"message_id"`
	Duration  *int `json:"duration"`
}

type ArchiveItem struct {
	HistoryKey  string                  `json:"history_key"`
	Description *ArchiveItemDescription `json:"desc,omitempty"`
	Meta        *ArchiveItemMeta        `json:"meta,omitempty"`
	Chan        *ArchiveItemChan        `json:"chan,omitempty"`
}

type ArchiveItemDescription struct {
	Date      string   `json:"date"`
	Source    string   `json:"source"`
	ShowId    string   `json:"show_id"`
	ShowTitle string   `json:"show_title"`
	EpTitle   string   `json:"ep_title"`
	Artists   []string `json:"artists"`
}

type ArchiveItemMeta struct {
	Filename string `json:"filename"`
	Duration *int   `json:"duration,omitempty"`
	Size     *int   `json:"size,omitempty"`
}

type ArchiveItemChan struct {
	MessageId int `json:"msg_id,omitempty"`
}
