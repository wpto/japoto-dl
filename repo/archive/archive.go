package archive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pgeowng/japoto-dl/model"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

var (
	ErrDuplicate     = errors.New("record already exists")
	ErrAlreadyExists = fmt.Errorf("record already exists")
)

type Interface interface {
	Migrate() error
	IsLoaded(key string) (status ArchiveEntryStatus, err error)
	Create(key string, status ArchiveEntryStatus, archiveItem model.ArchiveItem) (err error)
	SetStatus(key string, status ArchiveEntryStatus) (err error)
}

type Archive interface {
	Migrate(archive *sqlitex.Pool) error
	IsLoaded(archive *sqlitex.Pool, key string) (ok bool, err error)
	SetLoaded(archive *sqlitex.Pool, key string, status bool) error
	Create(archive *sqlitex.Pool, key string, status bool, archiveItem model.ArchiveItem) error
}

func CreateDB(filename string) (pool *sqlitex.Pool, err error) {
	pool, err = sqlitex.Open(filename, sqlite.SQLITE_OPEN_CREATE|sqlite.SQLITE_OPEN_READWRITE, 10)
	if err != nil {
		err = fmt.Errorf("CreateDB: sqlitex.Open: %w", err)
		return
	}

	return
}

type ArchiveRepo struct {
	pool *sqlitex.Pool
}

func NewRepo(pool *sqlitex.Pool) (archive *ArchiveRepo, err error) {
	return &ArchiveRepo{
		pool: pool,
	}, nil
}

func (r *ArchiveRepo) Exec(ctx context.Context, query string) (err error) {
	conn := r.pool.Get(ctx)
	if conn == nil {
		return fmt.Errorf("failed to get connection")
	}
	defer r.pool.Put(conn)

	stmt := conn.Prep(query)
	defer func() {
		if err == nil {
			if cerr := stmt.Finalize(); cerr != nil {
				err = fmt.Errorf("ArchiveRepo.Exec: Finalize: %w", cerr)
				return
			}
		}
	}()

	if _, err = stmt.Step(); err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) Migrate() (err error) {
	conn := r.pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer r.pool.Put(conn)

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

func (r *ArchiveRepo) IsLoaded(key string) (status ArchiveEntryStatus, err error) {
	conn := r.pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer r.pool.Put(conn)

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

func (r *ArchiveRepo) SetStatus(key string, status ArchiveEntryStatus) (err error) {
	conn := r.pool.Get(context.TODO())
	if conn == nil {
		return
	}
	defer r.pool.Put(conn)

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

func (r *ArchiveRepo) Create(key string, status ArchiveEntryStatus, archiveItem model.ArchiveItem) (err error) {
	bytes, err := json.Marshal(archiveItem)
	if err != nil {
		return fmt.Errorf("ArchiveRepo.Create: can't marshal archiveItem: %w", err)
	}

	conn := r.pool.Get(context.TODO())
	if conn == nil {
		return fmt.Errorf("ArchiveRepo.Create: couldn't get connection for pool")
	}
	defer r.pool.Put(conn)

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

func (r *ArchiveRepo) InsertShow(ctx context.Context, source string, showName string) (err error) {
	conn := r.pool.Get(ctx)
	if conn == nil {
		return fmt.Errorf("archive.InsertShow: couldn't get connection")
	}
	defer r.pool.Put(conn)

	query := `insert into shows(source, show_name, inserted_at, updated_at) values ($source, $show_name, datetime('now'), '2020-01-01 00:00:00') on conflict(source, show_name) do nothing;`
	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetText("$source", source)
	stmt.SetText("$show_name", showName)

	_, err = stmt.Step()

	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return
}

func (r *ArchiveRepo) LockShowsForUpdate(ctx context.Context, workerID string, source string) (err error) {
	conn := r.pool.Get(ctx)
	if conn == nil {
		return fmt.Errorf("archive.LockShowsForUpdate: couldn't get connection")
	}
	defer r.pool.Put(conn)

	// update shows set worker_id = 'fqwfg', updated_at = datetime('now', '+1 hour') from ( select source, show_name from shows where source = 'hibiki' and datetime(updated_at, '+1 hour') < datetime('now') limit 5 ) as src where shows.source = src.source and shows.show_name = src.show_name;
	query := `update shows set worker_id = $worker_id,
	updated_at = datetime('now', '+1 hour'),
	state = 'ready'
	where show_name in (
		select show_name from shows
		where source = $source
		  and (
		  	(datetime(updated_at, '+1 hour') < datetime('now') and worker_id = '')
		    or (datetime(updated_at, '+2 hour') < datetime('now') and coalesce(worker_id, '') <> '')
		  )
		limit 5 )
		and shows.source = $source;`

	stmt := conn.Prep(query)
	defer func() {
		if err != nil {
			fmt.Println("archive.InsertShow failed: %s", err.Error())
			if cerr := stmt.Finalize(); cerr != nil {
				err = fmt.Errorf("archive.InsertShow: finalize: %w", cerr)
				return
			}
		}
	}()

	stmt.SetText("$worker_id", workerID)
	stmt.SetText("$source", source)

	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

type ShowRow struct {
	ID   int64
	Name string
}

func (r *ArchiveRepo) GetLockedShows(ctx context.Context, workerID string) (show ShowRow, err error) {
	conn := r.pool.Get(ctx)
	if conn == nil {
		return ShowRow{}, fmt.Errorf("archive.LockShowsForUpdate: couldn't get connection")
	}
	defer r.pool.Put(conn)

	stmt := conn.Prep("SELECT id, show_name FROM shows WHERE worker_id = $worker_id and state <> 'done' limit 1;")
	stmt.SetText("$worker_id", workerID)

	found := false
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return ShowRow{}, fmt.Errorf("archive.GetLockedShows failed: %w", err)
		} else if !hasRow {
			break
		} else if !found {
			show.ID = stmt.GetInt64("id")
			show.Name = stmt.GetText("show_name")
			found = true
		}
	}

	if !found {
		return ShowRow{}, fmt.Errorf("archive.GetLockedShows failed: noRows")
	}

	return
}

func (r *ArchiveRepo) UpdateLockedShow(ctx context.Context, showID int64, state string) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	q := `update shows set state = $state where id = $id`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetInt64("$id", showID)
	stmt.SetText("$state", state)

	if _, err := stmt.Step(); err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) GetConn(ctx context.Context) (conn *sqlite.Conn, put func(), err error) {
	conn = r.pool.Get(ctx)
	if conn == nil {
		return nil, nil, fmt.Errorf("GetConn failed")
	}
	return conn, func() { r.pool.Put(conn) }, nil
}

func (r *ArchiveRepo) Finalize(stmt *sqlite.Stmt) {
	if err := stmt.Finalize(); err != nil {
		fmt.Println("r.Finalize failed: %s", err.Error())
		return
	}
}

func (r *ArchiveRepo) EnsurePersona(ctx context.Context, name string) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	query := `insert into persona(name, inserted_at) values ($name, datetime('now')) on conflict(name) do nothing`
	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetText("$name", name)

	_, err = stmt.Step()
	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("step failed: %w", err)
	}

	return nil
}

func (r *ArchiveRepo) GetPersona(ctx context.Context, name string) (int64, error) {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return 0, err
	}
	defer finish()

	query := `select id from persona where name = $name order by id limit 1`
	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetText("$name", name)

	if hasRow, err := stmt.Step(); err != nil || !hasRow {
		return 0, fmt.Errorf("step failed or no rows: %w", err)
	}

	return stmt.GetInt64("id"), nil
}

func (r *ArchiveRepo) EnsurePersonaRole(ctx context.Context, roleName string, showID int64, personaID int64) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	query := `insert into persona_role(name, show_id, persona_id, inserted_at)
	values ($name, $show_id, $persona_id, datetime('now'))
	on conflict do nothing;`
	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetText("$name", roleName)
	stmt.SetInt64("$show_id", showID)
	stmt.SetInt64("$persona_id", personaID)

	_, err = stmt.Step()
	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) EnsurePersonaShowRelation(ctx context.Context, showID int64, personaID int64) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	query := `insert into persona_show(show_id, persona_id, inserted_at)
	values ($show_id, $persona_id, datetime('now'))
	on conflict do nothing;`
	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetInt64("$show_id", showID)
	stmt.SetInt64("$persona_id", personaID)

	_, err = stmt.Step()
	if err := isErrAlreadyExists(err); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) EnsureLatestKeyValue(ctx context.Context, showID int64, keyLabel, keyValue string) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	query := `with 
		key_ids as (select id from show_keys where label = $key_label)
		,new_values as (
			select column1 as show_id,
				column2 as title from (values ($show_id, $key_value))
		)
		,old_values as (
			select show_id, value as title
			from show_values
			where key in (select id from key_ids)
			  and show_id in (select show_id from new_values)
			order by inserted_at desc
			limit 1
		)
		insert into show_values(show_id, key, value, inserted_at)
		select
		  new_values.show_id,
		  key_ids.id as key,
		  new_values.title,
		  datetime('now') as inserted_at
		from new_values
		left join key_ids 
		left join old_values
    where old_values.title is null or old_values.title <> new_values.title`

	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetInt64("$show_id", showID)
	stmt.SetText("$key_label", keyLabel)
	stmt.SetText("$key_value", keyValue)

	if _, err := stmt.Step(); err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) EnsureURLBank(ctx context.Context, url string) (int64, error) {
	url = strings.TrimSpace(url)

	if err := r.saveURLToBank(ctx, url); err != nil && !errors.Is(err, ErrAlreadyExists) {
		return 0, fmt.Errorf("saveURLToBank failed: %w", err)
	}

	key, err := r.getURLFromBank(ctx, url)
	if err != nil {
		return 0, fmt.Errorf("getURLFromBank failed: %w", err)
	}

	return key, nil
}

func (r *ArchiveRepo) saveURLToBank(ctx context.Context, url string) (err error) {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	q := `insert into url_bank (url, inserted_at) values ($url, datetime('now'))`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetText("$url", url)

	_, err = stmt.Step()
	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) getURLFromBank(ctx context.Context, url string) (key int64, err error) {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return 0, err
	}
	defer finish()

	q := `select id from url_bank where url = $url limit 1`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetText("$url", url)

	hasRow, err := stmt.Step()
	if err != nil {
		return 0, err
	}

	if !hasRow {
		return 0, fmt.Errorf("no rows")
	}

	return stmt.GetInt64("id"), nil
}

func (r *ArchiveRepo) EnsureEpisode(ctx context.Context, ep AEpisode) (int64, error) {
	err := r.insertEpisode(ctx, ep)
	if err != nil {
		return 0, fmt.Errorf("insertEpisode failed: %w", err)
	}

	id, err := r.getEpisodeID(ctx, ep)
	if err != nil {
		return 0, fmt.Errorf("getEpisodeID failed: %w", err)
	}

	return id, nil
}

type AEpisode struct {
	ShowID int64
	Title  string
	Date   ADate
}

type ADate struct {
	Day   int64
	Month int64
	Year  int64
}

func (r *ArchiveRepo) insertEpisode(ctx context.Context, ep AEpisode) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	q := `insert into show_episodes (
  	      show_id, title, date_day, date_month, date_year, inserted_at)
        values ($show_id, $title, $date_day, $date_month, $date_year, datetime('now'))
        on conflict do nothing`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetInt64("$show_id", ep.ShowID)
	stmt.SetText("$title", ep.Title)
	stmt.SetInt64("$date_day", ep.Date.Day)
	stmt.SetInt64("$date_month", ep.Date.Month)
	stmt.SetInt64("$date_year", ep.Date.Year)

	_, err = stmt.Step()

	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepo) getEpisodeID(ctx context.Context, ep AEpisode) (int64, error) {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return 0, err
	}
	defer finish()

	q := `select id from show_episodes
        where show_id = $show_id
        and title = $title
        order by inserted_at desc
        limit 1`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetInt64("$show_id", ep.ShowID)
	stmt.SetText("$title", ep.Title)

	hasRow, err := stmt.Step()
	if err != nil {
		return 0, err
	}

	if !hasRow {
		return 0, fmt.Errorf("no rows")
	}

	return stmt.GetInt64("id"), nil
}

func (r *ArchiveRepo) AddDownloadJob(ctx context.Context, job model.DownloadJob) error {
	fmt.Println("add download job", job)
	if err := r.saveURLToBank(ctx, job.PlaylistURL); err != nil && !errors.Is(err, ErrAlreadyExists) {
		return err
	}

	if err := r.saveURLToBank(ctx, job.ImageURL); err != nil && !errors.Is(err, ErrAlreadyExists) {
		return err
	}

	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	q := `insert into download_jobs
        (ep_id, playlist_url, image_url, inserted_at, worker_id, updated_at, status)
        select 
          column1 as ep_id,
          pu.id as playlist_url,
          iu.id as image_url,
          datetime('now') as inserted_at,
          '' as worker_id,
          datetime('now') as updated_at,
          'ready' as status
        from (values ($ep_id))
        left join url_bank as pu on pu.url = $playlist_url
        left join url_bank as iu on iu.url = $image_url`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetInt64("$ep_id", job.EpisodeID)
	stmt.SetText("$playlist_url", job.PlaylistURL)
	stmt.SetText("$image_url", job.ImageURL)

	_, err = stmt.Step()

	if err := isErrAlreadyExists(err); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func isErrAlreadyExists(err error) error {
	var serr sqlite.Error

	if err != nil && errors.As(err, &serr) && serr.Code == sqlite.SQLITE_CONSTRAINT_UNIQUE {
		return ErrAlreadyExists
	}

	return err
}

func (r *ArchiveRepo) LockDownloadJobs(ctx context.Context, workerID string, source string) (err error) {
	conn := r.pool.Get(ctx)
	if conn == nil {
		return fmt.Errorf("archive.LockDownloadJobs: couldn't get connection")
	}
	defer r.pool.Put(conn)

	// update shows set worker_id = 'fqwfg', updated_at = datetime('now', '+1 hour') from ( select source, show_name from shows where source = 'hibiki' and datetime(updated_at, '+1 hour') < datetime('now') limit 5 ) as src where shows.source = src.source and shows.show_name = src.show_name;
	query := `
	update download_jobs set
	  worker_id = $worker_id,
	  updated_at = datetime('now', '+1 hour'),
	  status = 'ready'
	where id in (
		select id from download_jobs
		where
			(status = '' or status = 'ready')
		  and (
		  	(datetime(updated_at, '+1 hour') < datetime('now') and worker_id = '')
		    or (datetime(updated_at, '+2 hour') < datetime('now') and coalesce(worker_id, '') <> '')
		  )
		limit 5
	);`

	stmt := conn.Prep(query)
	defer r.Finalize(stmt)

	stmt.SetText("$worker_id", workerID)

	if _, err = stmt.Step(); err != nil {
		return
	}

	return
}

type EpisodeDate struct {
	Year  int
	Month int
	Day   int
}

func (d EpisodeDate) String() string {
	result := ""

	if d.Year < 0 {
		result += "00"
	} else {
		result += fmt.Sprintf("%02d", d.Year%100)
	}

	result += fmt.Sprintf("%02d%02d", d.Month, d.Day)

	return result
}

type EpisodeLocalIdx int64

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func (e EpisodeLocalIdx) String() string {
	i := int64(e)
	ll := int64(len(base58Alphabet))
	result := ""
	for i > 0 {
		c := i % ll
		result += base58Alphabet[c : c+1]
		i = i / ll
	}
	if result == "" {
		return base58Alphabet[0:1]
	}
	return result
}

func (e EpisodeLocalIdx) i64() int64 {
	return int64(e)
}

type DownloadJob struct {
	JobID       int64
	PlaylistURL string
	ImageURL    string

	ShowSource string
	ShowName   string
	ShowTitle  string

	EpisodeDate     EpisodeDate
	EpisodeTitle    string
	EpisodeLocalIdx EpisodeLocalIdx
}

func (r *ArchiveRepo) GetLockedDownloadJob(ctx context.Context, workerID string) (DownloadJob, error) {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return DownloadJob{}, err
	}
	defer finish()

	stmt := conn.Prep(`
		SELECT
		  jobs.id as job_id,
		  ub1.url as playlist_url,
		  ub2.url as image_url,
      sh.source as source,
      se.date_year as date_year,
      se.date_month as date_month,
      se.date_day as date_day,
      sh.show_name as show_name,
      se.title as episode_title,
      sv.value as show_title,
      se.id as episode_local_idx

		FROM download_jobs jobs
		LEFT JOIN url_bank ub1 ON ub1.id = jobs.playlist_url
		LEFT JOIN url_bank ub2 ON ub2.id = jobs.image_url
		LEFT JOIN show_episodes se ON se.id = jobs.ep_id
		LEFT JOIN shows sh ON sh.id = se.show_id

		LEFT JOIN show_values sv ON sv.show_id = sh.id
		LEFT JOIN show_keys sk ON sk.id = sv.key AND sk.label = 'title'

		WHERE
		  jobs.worker_id = $worker_id
		  and jobs.status = 'ready'
		LIMIT 1;`)

	stmt.SetText("$worker_id", workerID)
	defer r.Finalize(stmt)

	var job DownloadJob
	found := false
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return job, fmt.Errorf("archive.GetLockedDownloadJob failed: %w", err)
		} else if !hasRow {
			break
		} else if !found {
			job.JobID = stmt.GetInt64("job_id")

			job.PlaylistURL = stmt.GetText("playlist_url")
			job.ImageURL = stmt.GetText("image_url")

			job.ShowSource = stmt.GetText("show_source")
			job.ShowName = stmt.GetText("show_name")
			job.ShowTitle = stmt.GetText("show_title")

			job.EpisodeDate.Year = int(stmt.GetInt64("date_year"))
			job.EpisodeDate.Month = int(stmt.GetInt64("date_month"))
			job.EpisodeDate.Day = int(stmt.GetInt64("date_day"))
			job.EpisodeTitle = stmt.GetText("episode_title")
			job.EpisodeLocalIdx = EpisodeLocalIdx(stmt.GetInt64("episode_local_idx"))

			found = true
		}
	}

	if !found {
		return job, fmt.Errorf("archive.GetLockedDownloadJob failed: noRows")
	}

	return job, nil
}

type Personas []model.Performer

func (p Personas) String() string {
	names := make([]string, 0, 2)
	for _, pp := range p {
		name := strings.TrimSpace(pp.Name)
		if name != "" {
			names = append(names, name)
		}
		role := strings.TrimSpace(pp.Role)
		if role != "" {
			names = append(names, role)
		}
	}

	return strings.Join(names, " ")
}

func (r *ArchiveRepo) GetEpisodePersona(ctx context.Context, episodeLocalIdx EpisodeLocalIdx) (Personas, error) {
	fmt.Println("GetEpisodePersona: ", episodeLocalIdx.i64())
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer finish()

	stmt := conn.Prep(`
		SELECT
			p.name as persona_name,
			coalesce(pr.name, '') as persona_role

		FROM show_episodes se
		LEFT JOIN persona_show ps ON ps.show_id = se.show_id
		LEFT JOIN persona p ON p.id = ps.persona_id
		LEFT JOIN persona_role pr
		  ON pr.show_id = se.show_id
		  AND pr.persona_id = p.id
		WHERE
      se.id = $episode_local_idx`)

	stmt.SetInt64("$episode_local_idx", episodeLocalIdx.i64())
	defer r.Finalize(stmt)

	persona := make([]model.Performer, 0, 2)
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("archive.GetEpisodePersona failed: %w", err)
		} else if !hasRow {
			break
		}

		var p model.Performer

		p.Name = stmt.GetText("persona_name")
		p.Role = stmt.GetText("persona_role")

		persona = append(persona, p)
	}

	return persona, nil
}

func (r *ArchiveRepo) UpdateLockedDownloadJob(ctx context.Context, jobID int64, state string) error {
	conn, finish, err := r.GetConn(ctx)
	if err != nil {
		return err
	}
	defer finish()

	q := `update download_jobs set status = $state where id = $id`
	stmt := conn.Prep(q)
	defer r.Finalize(stmt)

	stmt.SetInt64("$id", jobID)
	stmt.SetText("$state", state)

	if _, err := stmt.Step(); err != nil {
		return err
	}

	return nil
}
