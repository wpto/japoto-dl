package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pgeowng/japoto-dl/pkg/dl"
	"github.com/pgeowng/japoto-dl/provider"
	"github.com/pgeowng/japoto-dl/tasks"
	"github.com/pgeowng/japoto-dl/workdir"
	"github.com/pgeowng/japoto-dl/workdir/muxer"
	"github.com/pgeowng/japoto-dl/workdir/wd"

	"github.com/pgeowng/japoto-dl/repo/archive"
	status1 "github.com/pgeowng/japoto-dl/repo/status"
)

type PostDownloadHook = func(wdhls *workdir.Workdir, wd *wd.Wd)

type DownloadWorker struct {
	repo             *archive.ArchiveRepo
	postDownloadHook PostDownloadHook
}

func NewDownloadWorker(repo *archive.ArchiveRepo) *DownloadWorker {
	return &DownloadWorker{
		repo: repo,
	}
}

func (a *DownloadWorker) WithPostDownloadHook(hook PostDownloadHook) *DownloadWorker {
	a.postDownloadHook = hook
	return a
}

func (a *DownloadWorker) RunDownloadWorker(ctx context.Context) error {
	log.Printf("Download worker is running\n")

	workerID := NewWorkerID()
	workerID = "1"

	if err := a.repo.LockDownloadJobs(ctx, workerID, "some source"); err != nil {
		fmt.Println(err)
		return err
	}

	job, err := a.repo.GetLockedDownloadJob(ctx, workerID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	personas, err := a.repo.GetEpisodePersona(ctx, job.EpisodeLocalIdx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	tags := map[string]string{
		"title":  strings.Join([]string{job.EpisodeDate.String(), job.ShowName, job.EpisodeTitle, job.ShowTitle}, " "),
		"artist": personas.String(),
		"album":  job.ShowTitle,
		"track":  job.EpisodeDate.String(),
	}

	salt := fmt.Sprintf("%s-%s--%s-u%s", job.EpisodeDate.String(), job.ShowName, job.ShowSource, job.EpisodeLocalIdx.String())
	filename := fmt.Sprintf("%s.mp3", salt)

	destPath := fmt.Sprintf("./%s", filename)

	ffm := muxer.NewFFMpegHLS(destPath, tags)
	wd1 := wd.NewWd("./.cache", salt)

	wd := workdir.NewWorkdir(wd1, ffm, map[string]string{
		"playlist": "playlist.m3u8",
		"image":    "image",
	})
	hls := tasks.NewTasks(wd)

	fmt.Println(job)
	loader := dl.NewGrequests()
	metric := status1.NewMetric()

	genLoader := provider.NewGeneralLoader(loader, hls.AudioHLS(), metric, wd)

	err = genLoader.DownloadEpisode(ctx, provider.DownloadEpisodeParams{
		PlaylistURL:    job.PlaylistURL,
		ImageURL:       job.ImageURL,
		RequestOptions: provider.OnsenGopts,
	})
	if err != nil {
		fmt.Println(err)
		if err := a.repo.UpdateLockedDownloadJob(ctx, job.JobID, "failed"); err != nil {
			fmt.Println("UpdateLockedDownloadJob: failed: %w", err)
		}
		return err
	}

	if err := a.repo.UpdateLockedDownloadJob(ctx, job.JobID, "done"); err != nil {
		fmt.Println("UpdateLockedDownloadJob: failed: %w", err)
	}

	return nil
}
