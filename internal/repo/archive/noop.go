package archive

// type ReadOnlyArchiveRepo struct{ Archive }

// func (r *ReadOnlyArchiveRepo) SetLoaded(key string, loaded bool) (err error)          { return nil }
// func (r *ReadOnlyArchiveRepo) Create(key string, archiveItem ArchiveItem) (err error) { return nil }
// func NewReadOnlyRepo(filename string) (archive Archive, err error) {
// 	archiveRepo, err := NewRepo(filename)
// 	if err != nil {
// 		return
// 	}
// 	archive = &WriteOnlyArchiveRepo{Archive: archiveRepo}
// 	return
// }

// type WriteOnlyArchiveRepo struct{ Archive }

// func (r *WriteOnlyArchiveRepo) IsLoaded(key string) (ok bool, err error) { return false, nil }
// func NewWriteOnlyRepo(filename string) (archive Archive, err error) {
// 	archiveRepo, err := NewRepo(filename)
// 	if err != nil {
// 		return
// 	}
// 	archive = &WriteOnlyArchiveRepo{Archive: archiveRepo}
// 	return
// }

// type NoopArchiveRepo struct{}

// func (r *NoopArchiveRepo) IsLoaded(key string) (ok bool, err error)               { return false, nil }
// func (r *NoopArchiveRepo) SetLoaded(key string, loaded bool) (err error)          { return nil }
// func (r *NoopArchiveRepo) Create(key string, archiveItem ArchiveItem) (err error) { return nil }
// func NewNoopRepo() Archive                                                        { return &NoopArchiveRepo{} }
