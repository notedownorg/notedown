package workspace

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func (w *Workspace) DailyNotePath(date time.Time) (string, error) {
	path := filepath.Join(w.directories.DailyNotes, date.Format("2006-01-02.md"))

	// Create daily note if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return "", err
		}
		file.Close()
		w.cache.Wait(path, time.Second) // Wait for the cache to be updated before returning
		slog.Info("created daily note", slog.String("path", path))
	}

	return path, nil
}
