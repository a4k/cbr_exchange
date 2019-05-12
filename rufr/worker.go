package rufr

import (
	"context"
	"net/http"
	"time"

	"../app"
	"../jobs"
	l4g "../utils/log4go"
	"../model"
)

const (
	TIME_BETWEEN_BATCHES = 30
)

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	app       *app.App
}

func (m *FinHelpJobInterfaceImpl) MakeWorker() model.Worker {
	worker := Worker{
		name:      "FinHelp",
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: m.App.Jobs,
		app:       m.App,
	}

	return &worker
}

func (worker *Worker) Run() {
	l4g.Debug("РУФР начал работу %s",worker.name)

	defer func() {
		l4g.Debug("РУФР закончил работу %s %s", "worker", worker.name)
		worker.stopped <- true
	}()

	for {
		select {
		case <-worker.stop:
			l4g.Debug("Скрипт получил стоп сигнал %s %s", "worker", worker.name)
			return
		case job := <-worker.jobs:
			l4g.Debug("Скрипт получил нового кандидата. %s %s", "worker", worker.name)
			worker.DoJob(&job)
		}
	}
}

func (worker *Worker) Stop() {
	l4g.Debug("Остановка скрипта %s %s", "worker", worker.name)
	worker.stop <- true
	<-worker.stopped
}

func (worker *Worker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		l4g.Info("У работника возникла ошибка при попытке начала работы",
			"worker", worker.name,
			"job_id", job.Id,
			"error", err.Error())
		return
	} else if !claimed {
		return
	}

	cfg := worker.app.Config()

	timeoutSeconds := time.Duration(*cfg.ExchangeSettings.RequestTimeoutSeconds)

	cancelCtx, cancelCancelWatcher := context.WithCancel(context.Background())
	cancelWatcherChan := make(chan interface{}, 1)
	go worker.app.Jobs.CancellationWatcher(cancelCtx, job.Id, cancelWatcherChan)

	defer cancelCancelWatcher()

	for {
		select {
		case <-cancelWatcherChan:
			l4g.Debug("Работа отменена через CancellationWatcher %s %s %s %s", "worker", worker.name, "job_id", job.Id)
			worker.setJobCanceled(job)
			return

		case <-worker.stop:
			l4g.Debug("Работа отменена через Worker Stop %s %s %s %s", "worker", worker.name, "job_id", job.Id)
			worker.setJobCanceled(job)
			return

		case <-time.After(timeoutSeconds * time.Second):
			done, progress, err := worker.runFinHelp(job.Data[JOB_DATA_KEY_FINHELP], job.Data[JOB_DATA_KEY_FINHELP_LAST_DONE])
			if err != nil {
				l4g.Error("Не удалось запустить скрипт", "worker", worker.name, "job_id", job.Id, "error", err.Error())
				worker.setJobError(job, err)
				return
			} else if done {
				l4g.Info("Работа завершена %s %s %s %s", "worker", worker.name, "job_id", job.Id)
				worker.setJobSuccess(job)
				return
			} else {
				job.Data[JOB_DATA_KEY_FINHELP_LAST_DONE] = progress
				if err := worker.app.Jobs.UpdateInProgressJobData(job); err != nil {
					l4g.Error("Не удалось обновить статус для работы", "worker", worker.name, "job_id", job.Id, "error", err.Error())
					worker.setJobError(job, err)
					return
				}
			}
		}
	}
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.app.Jobs.SetJobSuccess(job); err != nil {
		l4g.Error("Не удалось задать успех для работы %s %s %s %s", "worker", worker.name, "job_id", job.Id, "error", err.Error())
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.app.Jobs.SetJobError(job, appError); err != nil {
		l4g.Error("Не удалось задать ошибку для работы %s %s %s %s", "worker", worker.name, "job_id", job.Id, "error", err.Error())
	}
}

func (worker *Worker) setJobCanceled(job *model.Job) {
	if err := worker.app.Jobs.SetJobCanceled(job); err != nil {
		l4g.Error("Не удалось отметить задание отмененным %s %s %s %s", "worker", worker.name, "job_id", job.Id, "error", err.Error())
	}
}

// Return parameters:
// - whether the finhelp is completed on this run (true) or still incomplete (false).
// - the updated lastDone string for the finhelp.
// - any error which may have occurred while running the finhelp.
func (worker *Worker) runFinHelp(key string, lastDone string) (bool, string, *model.AppError) {
	var done bool
	var progress string
	var err *model.AppError

	switch key {
	case model.FINHELP_KEY_ADVANCED_PARSING:
		done, progress, err = worker.runAdvancedParsing(lastDone)
	default:
		return false, "", model.NewAppError("FinHelpWorker.runFinHelp", "finhelp.worker.run_finhelp.unknown_key",
			map[string]interface{}{"key": key}, "", http.StatusInternalServerError)
	}

	return done, progress, err
}
