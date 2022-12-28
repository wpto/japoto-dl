package jobstorage

import "crawshaw.io/sqlite/sqlitex"

type JobStorage struct {
	pool *sqlitex.Pool
	log  *zerolog.Log
}

func NewJobStorage() (jobStorage *JobStorage, err error) {
	return &JobStorage{}, nil
}
