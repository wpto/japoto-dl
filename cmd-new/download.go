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

	status1 "github.com/pgeowng/japoto-dl/repo/status"
)

func (a *App) RunDownloadWorker(ctx context.Context) error {
	log.Printf("Download worker is running\n")

	workerID := NewWorkerID()
	workerID = "1"

	if err := a.arch.LockDownloadJobs(ctx, workerID, "some source"); err != nil {
		fmt.Println(err)
		return err
	}

	job, err := a.arch.GetLockedDownloadJob(ctx, workerID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	personas, err := a.arch.GetEpisodePersona(ctx, job.EpisodeLocalIdx)
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

	err = genLoader.DownloadEpisode(provider.DownloadEpisodeParams{
		PlaylistURL:    job.PlaylistURL,
		ImageURL:       job.ImageURL,
		RequestOptions: provider.HibikiGopts,
	})
	if err != nil {
		fmt.Println(err)
		if err := a.arch.UpdateLockedDownloadJob(ctx, job.JobID, "failed"); err != nil {
			fmt.Println("UpdateLockedDownloadJob: failed: %w", err)
		}
		return err
	}

	if err := a.arch.UpdateLockedDownloadJob(ctx, job.JobID, "done"); err != nil {
		fmt.Println("UpdateLockedDownloadJob: failed: %w", err)
	}

	return nil
}
