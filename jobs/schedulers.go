package jobs

import (
	"sync"
	"time"

	l4g "../utils/log4go"

	"../model"
)

type Schedulers struct {
	stop          chan bool
	stopped       chan bool
	configChanged chan *model.Config
	listenerId    string
	startOnce     sync.Once
	jobs          *JobServer

	schedulers   []model.Scheduler
	nextRunTimes []*time.Time
}

func (srv *JobServer) InitSchedulers() *Schedulers {
	l4g.Debug("Инициализация планировщиков.")

	schedulers := &Schedulers{
		stop:          make(chan bool),
		stopped:       make(chan bool),
		configChanged: make(chan *model.Config),
		jobs:          srv,
	}

	if migrationsInterface := srv.Migrations; migrationsInterface != nil {
		schedulers.schedulers = append(schedulers.schedulers, migrationsInterface.MakeScheduler())
	}
	if finhelpInterface := srv.FinHelp; finhelpInterface != nil {
		schedulers.schedulers = append(schedulers.schedulers, finhelpInterface.MakeScheduler())
	}

	schedulers.nextRunTimes = make([]*time.Time, len(schedulers.schedulers))
	return schedulers
}


func (schedulers *Schedulers) Start() *Schedulers {
	schedulers.listenerId = schedulers.jobs.ConfigService.AddConfigListener(schedulers.handleConfigChange)

	go func() {
		schedulers.startOnce.Do(func() {
			l4g.Info("Запуск планировщиков.")

			defer func() {
				l4g.Info("Остановка планировщиков.")
				close(schedulers.stopped)
			}()

			now := time.Now()
			for idx, scheduler := range schedulers.schedulers {
				if !scheduler.Enabled(schedulers.jobs.Config()) {
					schedulers.nextRunTimes[idx] = nil
				} else {
					schedulers.setNextRunTime(schedulers.jobs.Config(), idx, now, false)
				}
			}

			for {
				select {
				case <-schedulers.stop:
					l4g.Debug("Планировщик получил стоп-сигнал.")
					return
				case now = <-time.After(5 * time.Second):
					cfg := schedulers.jobs.Config()

					for idx, nextTime := range schedulers.nextRunTimes {
						if nextTime == nil {
							continue
						}

						if time.Now().After(*nextTime) {
							scheduler := schedulers.schedulers[idx]
							if scheduler != nil {
								if scheduler.Enabled(cfg) {
									if _, err := schedulers.scheduleJob(cfg, scheduler); err != nil {
										l4g.Warn("Не удалось запланировать работу: %v", scheduler.Name())
										l4g.Error(err)
									} else {
										schedulers.setNextRunTime(cfg, idx, now, true)
									}
								}
							}
						}
					}
				case newCfg := <-schedulers.configChanged:
					for idx, scheduler := range schedulers.schedulers {
						if !scheduler.Enabled(newCfg) {
							schedulers.nextRunTimes[idx] = nil
						} else {
							schedulers.setNextRunTime(newCfg, idx, now, false)
						}
					}
				}
			}
		})
	}()

	return schedulers
}

func (schedulers *Schedulers) Stop() *Schedulers {
	l4g.Info("Остановка планировщиков.")
	close(schedulers.stop)
	<-schedulers.stopped
	return schedulers
}

func (schedulers *Schedulers) setNextRunTime(cfg *model.Config, idx int, now time.Time, pendingJobs bool) {
	scheduler := schedulers.schedulers[idx]

	if !pendingJobs {
		if pj, err := schedulers.jobs.CheckForPendingJobsByType(scheduler.JobType()); err != nil {
			l4g.Error("Не удалось задать следующее время выполнения: %s", err.Error())
			schedulers.nextRunTimes[idx] = nil
			return
		} else {
			pendingJobs = pj
		}
	}

	lastSuccessfulJob, err := schedulers.jobs.GetLastSuccessfulJobByType(scheduler.JobType())
	if err != nil {
		l4g.Error("Не удалось задать следующее время выполнения: %s" + err.Error())
		schedulers.nextRunTimes[idx] = nil
		return
	}

	schedulers.nextRunTimes[idx] = scheduler.NextScheduleTime(cfg, now, pendingJobs, lastSuccessfulJob)
	l4g.Debug("Следующее время выполнения %v: %v", scheduler.Name(), schedulers.nextRunTimes[idx])
}

func (schedulers *Schedulers) scheduleJob(cfg *model.Config, scheduler model.Scheduler) (*model.Job, *model.AppError) {
	pendingJobs, err := schedulers.jobs.CheckForPendingJobsByType(scheduler.JobType())
	if err != nil {
		return nil, err
	}

	lastSuccessfulJob, err2 := schedulers.jobs.GetLastSuccessfulJobByType(scheduler.JobType())
	if err2 != nil {
		return nil, err
	}

	return scheduler.ScheduleJob(cfg, pendingJobs, lastSuccessfulJob)
}

func (schedulers *Schedulers) handleConfigChange(oldConfig *model.Config, newConfig *model.Config) {
	l4g.Debug("Получено изменение конфига.")
	schedulers.configChanged <- newConfig
}
