package rufr

import (
	"time"
	"../app"
	l4g "../utils/log4go"
	"../model"
	"../store"
)

type Scheduler struct {
	App                 *app.App
	allFinHelpCompleted bool
}

func (m *FinHelpJobInterfaceImpl) MakeScheduler() model.Scheduler {
	return &Scheduler{m.App, false}
}

func (scheduler *Scheduler) Name() string {
	return "РУФР ПЛАНИРОВЩИК"
}

func (scheduler *Scheduler) JobType() string {
	return model.JOB_TYPE_FINHELP
}

func (scheduler *Scheduler) Enabled(cfg *model.Config) bool {
	return true
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	if scheduler.allFinHelpCompleted {
		return nil
	}

	// Время вызова задания по расписанию
	timeoutSeconds := time.Duration(*cfg.ExchangeSettings.RequestTimeoutSeconds)
	nextTime := time.Now().Add(timeoutSeconds * time.Second)

	return &nextTime
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	l4g.Debug("Планирование задания %s", scheduler.Name())

	// Work through the list of finhelp in order. Schedule the first one that isn't done (assuming it isn't in progress already).
	for _, key := range MakeFinHelpList() {
		state, job, err := GetFinHelpState(key, scheduler.App.Srv.Store)
		if err != nil {
			l4g.Error("Не удалось определить статус: %s %s %s", scheduler.Name(), key, err.Error())
			return nil, nil
		}

		if state == FINHELP_STATE_IN_PROGRESS {
			// Check the finhelp job isn't wedged.
			return scheduler.createJob(key, job, scheduler.App.Srv.Store)
		}

		if state == FINHELP_STATE_COMPLETED {
			// This finhelp is done. Continue to check the next.
			continue
		}

		if state == FINHELP_STATE_UNSCHEDULED {
			l4g.Debug("Планирование нового задания для РУФР. %s %s %s %s", "scheduler", scheduler.Name(), "finhelp_key", key)
			return scheduler.createJob(key, job, scheduler.App.Srv.Store)
		}

		l4g.Error("Не известное состояние. %s %s", "finhelp_state", state)
		return nil, nil
	}

	// If we reached here, then there aren't any finhelp left to run.
	scheduler.allFinHelpCompleted = true
	l4g.Debug("Все планировщики завершены. %s %s", "scheduler", scheduler.Name())

	return nil, nil
}

func (scheduler *Scheduler) createJob(finhelpKey string, lastJob *model.Job, store store.Store) (*model.Job, *model.AppError) {
	var lastDone string
	if lastJob != nil {
		lastDone = lastJob.Data[JOB_DATA_KEY_FINHELP_LAST_DONE]
	}

	data := map[string]string{
		JOB_DATA_KEY_FINHELP:           finhelpKey,
		JOB_DATA_KEY_FINHELP_LAST_DONE: lastDone,
	}

	if job, err := scheduler.App.Jobs.CreateJob(model.JOB_TYPE_FINHELP, data); err != nil {
		return nil, err
	} else {
		return job, nil
	}
}
