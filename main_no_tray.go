//go:build no_tray

package main

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gen2brain/beeep"
	"github.com/k-x7/eventt"
	"golang.org/x/exp/slog"
)

const (
	addr   = "localhost"
	port   = "29172"
	dir    = "/events"
	title  = "Alertt"
	desc   = "monitor Sonarr events and show system notifications"
	poster = "Sonarr/MediaCover/%d/poster-250.jpg"
)

//go:embed res
var res embed.FS

var (
	iconPath string
)

var server *http.Server = &http.Server{
	Addr:    addr + ":" + port,
	Handler: nil,
}

func main() {
	setupIcon()
	setupSonarEvents()
}

func setupIcon() {
	// for system notification icon if poster not found locally
	b, err := res.ReadFile("res/gopher.png")
	if err != nil {
		slog.Error("can't read from embed fs", err)
	}

	// save system notification icon in temp dir
	tmp := os.TempDir()
	iconPath = filepath.Join(tmp, "/gopher.png")
	if err := os.WriteFile(iconPath, b, 0644); err != nil {
		slog.Error("can't save icon to temp dir", err)
	}
}

func setupSonarEvents() {

	// create sonarr events handler
	events := eventt.SonarrTriggers{
		// Log on errors
		LogOnError: true,

		// on grab show series title and session/episode number and release name
		// also include poster for the series if found
		OnGrab: func(event eventt.GrabEvent) {
			if err := beeep.Notify(
				fmt.Sprintf("Grabbed: '%s' S:%d, E:%d", event.Series.Title, event.Episodes[0].SeasonNumber, event.Episodes[0].EpisodeNumber),
				event.Release.ReleaseTitle,
				getPosterFromLocal(event.Series.ID),
			); err != nil {
				slog.Error("can't show notification for grab", err)
			}
		},

		// on download show series title and session/episode number and saved location
		// also include poster for the series if found
		// Note if there is any action
		OnDownload: func(event eventt.DownloadEvent) {
			if err := beeep.Notify(
				fmt.Sprintf("Downloaded: '%s' S:%d, E:%d", event.Series.Title, event.Episodes[0].SeasonNumber, event.Episodes[0].EpisodeNumber),
				event.EpisodeFile.Path,
				getPosterFromLocal(event.Series.ID),
			); err != nil {
				slog.Error("can't show notification for download", err)
			}
		},
	}

	slog.Info("start server", "host", addr, "port", port, "url", "http://"+addr+":"+port+dir)
	http.HandleFunc(dir, events.Monitor)
	if err := server.ListenAndServe(); err != nil {
		// if application quit don't show error
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error running http server", err)
		}
	}
}

// get poster path from local sonarr dir based on series id
func getPosterFromLocal(id int) string {
	posterPath := iconPath
	switch runtime.GOOS {
	case "windows":
		posterPath = filepath.Join("C:/ProgramData", fmt.Sprintf(poster, id))
	case "linux":
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			home = "~"
		}
		posterPath = filepath.Join(home, ".config", fmt.Sprintf(poster, id))
	}
	if exist, err := exists(posterPath); !exist || err != nil {
		posterPath = iconPath
	}
	return posterPath
}

// check if file exist or not with error if any
func exists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
