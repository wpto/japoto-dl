package main

import (
	"context"
	"fmt"

	"github.com/pgeowng/japoto-dl/repo/archive"
	status1 "github.com/pgeowng/japoto-dl/repo/status"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/wd"
)

type MuxJob struct {
	wdhls *workdir.Workdir
	wd    *wd.Wd
}

type MuxWorker struct {
	input  <-chan MuxJob
	done   chan struct{}
	status *status1.PrintLine
	repo   *archive.ArchiveRepo
}

func NewMuxWorker(status *status1.PrintLine, input <-chan MuxJob) *MuxWorker {
	ch := make(chan struct{}, 0)
	close(ch)
	return &MuxWorker{
		input:  input,
		done:   ch,
		status: status,
	}
}

func (a *MuxWorker) Done() <-chan struct{} {
	return a.done
}

func (a *MuxWorker) Start(ctx context.Context) {
	a.done = make(chan struct{}, 0)

	go func() {
		defer close(a.done)
		for job := range a.input {
			if err := a.iteration(ctx, job); err != nil {
				fmt.Println(err)
			}
		}
	}()
}

func (a *MuxWorker) iteration(ctx context.Context, job MuxJob) error {
	if err := job.wdhls.Mux(); err != nil {
		a.status.Error(fmt.Errorf("ffmpeg error: %w", err))
		return err
	}

	return nil
}
