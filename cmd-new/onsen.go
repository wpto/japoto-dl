package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pgeowng/japoto-dl/model"
	"github.com/pgeowng/japoto-dl/repo/archive"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewWorkerID() string {
	workerIDInt := int64(rand.Intn(12341245))
	workerID := strconv.FormatInt(workerIDInt, 36)
	return workerID
}

type Provider interface {
	GetShow(showName string) (model.Show, error)
	Label() string
}

func (a *App) RunGetRecentShows(ctx context.Context, provider Provider) (err error) {
	workerID := NewWorkerID()

	if err := a.arch.LockShowsForUpdate(ctx, workerID, provider.Label()); err != nil {
		fmt.Println(err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		showRow, err := a.arch.GetLockedShows(ctx, workerID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		log.Printf("Locked show: %#v\n", showRow)

		showEntity, err := provider.GetShow(showRow.Name)
		if err != nil {
			fmt.Println(err)
			return err
		}

		showTitle := strings.TrimSpace(showEntity.ShowTitle())
		err = a.arch.EnsureLatestKeyValue(ctx, showRow.ID, "title", showTitle)
		if err != nil {
			fmt.Println("write title failed:", err)
		}

		posterURL := strings.TrimSpace(showEntity.GeneralPosterURL())
		if err := a.writePosterURL(ctx, showRow.ID, posterURL); err != nil {
			fmt.Println(err)
			return err
		}

		performers := showEntity.GeneralPerformerInfo()
		fmt.Println(performers)
		if err := a.writePerformers(ctx, showRow.ID, performers); err != nil {
			fmt.Println(err)
			return err
		}

		genericEpisodes, err := showEntity.GetEpisodes(a.loader)
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, episode := range genericEpisodes {
			// onsenEp, ok := episode.(*provider.OnsenEpisode)
			// if !ok {
			// 	fmt.Println("onsen episode type assert failed")
			// 	return fmt.Errorf("onsen episode type assert failed")
			// }

			dd, mm, yy := episode.LeastDate()

			aep := archive.AEpisode{
				ShowID: showRow.ID,
				Title:  episode.EpTitle(),
				Date: archive.ADate{
					Day:   dd,
					Month: mm,
					Year:  yy,
				},
			}

			episodeID, err := a.arch.EnsureEpisode(ctx, aep)
			if err != nil && !errors.Is(err, archive.ErrAlreadyExists) {
				fmt.Println("ensureEpisode:", err)
				return err
			}

			jobs := episode.GetDownloadJobs(episodeID)
			for _, job := range jobs {
				err := a.arch.AddDownloadJob(ctx, job)
				if err != nil {
					fmt.Println(err)
				}
			}

			fmt.Println(episodeID)

		}

		err = a.arch.UpdateLockedShow(ctx, showRow.ID, "done")
		if err != nil {
			return err
		}

		fmt.Println(showRow, showEntity)
	}
	return nil
}

func (a *App) writePosterURL(ctx context.Context, showID int64, posterURL string) error {
	posterURLID, err := a.arch.EnsureURLBank(ctx, posterURL)
	if err != nil {
		return fmt.Errorf("writePosterURL:", err)
	}

	if len(posterURL) > 0 {
		err = a.arch.EnsureLatestKeyValue(ctx, showID, "poster", strconv.FormatInt(posterURLID, 10))
		if err != nil {
			return fmt.Errorf("writePosterURL: %w", err)
		}
	}

	return nil
}

func (a *App) writePerformers(ctx context.Context, showID int64, performers []model.Performer) error {
	for _, persona := range performers {
		err := a.arch.EnsurePersona(ctx, persona.Name)
		if err != nil {
			return fmt.Errorf("writePerformers: %w", err)
		}

		personaID, err := a.arch.GetPersona(ctx, persona.Name)
		if err != nil {
			return fmt.Errorf("writePerformers: %w", err)
		}

		if persona.Role != "" {
			err := a.arch.EnsurePersonaRole(ctx, persona.Role, showID, personaID)
			if err != nil {
				return fmt.Errorf("writePerformers: %w", err)
			}
		}

		err = a.arch.EnsurePersonaShowRelation(ctx, showID, personaID)
		if err != nil {
			return fmt.Errorf("writePerformers: %w", err)
		}
	}

	return nil
}
