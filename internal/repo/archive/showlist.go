package archive

import (
	"context"
	"fmt"
	"log"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

const memoryFile = "file:memory:?mode=memory"

type Storage struct {
	*sqlitex.Pool
}

func NewFileSQLite(filename string) (*sqlitex.Pool, error) {
	pool, err := sqlitex.Open(filename, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE, 10)
	if err != nil {
		return nil, fmt.Errorf("Failed to open memory sqlite: %w", err)
	}

	return pool, nil
}

// func NewFileSQLite(filename string) (*sqlitex.Pool, err)

func NewStorage(pool *sqlitex.Pool) *Storage {
	return &Storage{Pool: pool}
}

func (s *Storage) MigrateShowFeed(ctx context.Context) error {
	conn := s.Get(ctx)
	if conn == nil {
		return fmt.Errorf("Failed to get connection from pool")
	}

	defer s.Put(conn)

	query := `create table if not exists show_feed(
		source string not null,
		show_id string not null,
		feed_url string,
		last_update string,
		primary key (source, show_id)
	)`

	stmt, err := conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}

	defer func() {
		if err := stmt.Finalize(); err != nil {
			err = fmt.Errorf("Failed to finalize statement: %w", err)
		}
	}()

	if _, err := stmt.Step(); err != nil {
		return fmt.Errorf("Failed to migrate show feed: %w", err)
	}

	return nil
}

type ShowEntry struct {
	Source string
	ShowID string
}

func (s *Storage) AddShows(ctx context.Context, entries []ShowEntry) error {
	conn := s.Get(ctx)
	if conn == nil {
		return fmt.Errorf("Failed to get connection from pool")
	}

	defer s.Put(conn)

	query := `insert into show_feed (source, show_id) values ($source, $show_id) on conflict do nothing`

	stmt, err := conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}

	for _, entry := range entries {
		stmt.SetText("$source", entry.Source)
		stmt.SetText("$show_id", entry.ShowID)

		if _, err := stmt.Step(); err != nil {
			log.Default().Printf("Insert entry: %s %s", entry.Source, entry.ShowID)
			return fmt.Errorf("Failed to add show: %w", err)
		}

		stmt.Reset()
		stmt.ClearBindings()
	}

	if err := stmt.Finalize(); err != nil {
		return fmt.Errorf("Failed to finalize statement: %w", err)
	}

	return nil
}

func (s *Storage) GetShows(ctx context.Context) ([]ShowEntry, error) {
	conn := s.Get(ctx)
	if conn == nil {
		return nil, fmt.Errorf("Failed to get connection from pool")
	}

	defer s.Put(conn)

	result := make([]ShowEntry, 0, 16)

	query := `select source, show_id from show_feed`
	stmt := conn.Prep(query)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("Failed to read show_feed: %w", err)
		} else if !hasRow {
			break
		}

		result = append(result, ShowEntry{
			Source: stmt.GetText("source"),
			ShowID: stmt.GetText("show_id"),
		})
	}

	return result, nil
}
