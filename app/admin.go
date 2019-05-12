package app

import (
	"io"
	"os"
	"strings"
	"time"

	"net/http"

	l4g "../utils/log4go"
	"../model"
	"../utils"
)

func (a *App) GetLogs(page, perPage int) ([]string, *model.AppError) {
	var lines []string

	melines, err := a.GetLogsSkipSend(page, perPage)
	if err != nil {
		return nil, err
	}

	lines = append(lines, melines...)

	return lines, nil
}

func (a *App) GetLogsSkipSend(page, perPage int) ([]string, *model.AppError) {
	var lines []string

	if a.Config().LogSettings.EnableFile {
		file, err := os.Open(utils.GetLogFileLocation(a.Config().LogSettings.FileLocation))
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}

		defer file.Close()

		var newLine = []byte{'\n'}
		var lineCount int
		const searchPos = -1
		lineEndPos, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}
		for {
			pos, err := file.Seek(searchPos, io.SeekCurrent)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			b := make([]byte, 1)
			_, err = file.ReadAt(b, pos)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			if b[0] == newLine[0] || pos == 0 {
				lineCount++
				if lineCount > page*perPage {
					line := make([]byte, lineEndPos-pos)
					_, err := file.ReadAt(line, pos)
					if err != nil {
						return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
					}
					lines = append(lines, string(line))
				}
				if pos == 0 {
					break
				}
				lineEndPos = pos
			}

			if len(lines) == perPage {
				break
			}
		}

		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	} else {
		lines = append(lines, "")
	}

	return lines, nil
}

func (a *App) GetConfig() *model.Config {
	json := a.Config().ToJson()
	cfg := model.ConfigFromJson(strings.NewReader(json))
	cfg.Sanitize()

	return cfg
}

func (a *App) SaveConfig(cfg *model.Config, sendConfigChangeClusterMessage bool) *model.AppError {

	cfg.SetDefaults()
	a.Desanitize(cfg)

	if err := cfg.IsValid(); err != nil {
		return err
	}

	a.DisableConfigWatch()
	a.UpdateConfig(func(update *model.Config) {
		*update = *cfg
	})
	a.PersistConfig()
	a.ReloadConfig()
	a.EnableConfigWatch()

	return nil
}

func (a *App) RecycleDatabaseConnection() {
	oldStore := a.Srv.Store

	l4g.Warn(("api.admin.recycle_db_start.warn"))
	a.Srv.Store = a.newStore()
	a.Jobs.Store = a.Srv.Store

	if a.Srv.Store != oldStore {
		time.Sleep(20 * time.Second)
		oldStore.Close()
	}

	l4g.Warn(("api.admin.recycle_db_end.warn"))
}
