package jobs

import (
	"sync"

	l4g "../utils/log4go"
	"../model"
)

type Workers struct {
	startOnce     sync.Once
	ConfigService ConfigService
	Watcher       *Watcher
	Migrations    model.Worker
	listenerId    string
	FinHelp       model.Worker
}

func (srv *JobServer) InitWorkers() *Workers {
	workers := &Workers{
		ConfigService: srv.ConfigService,
	}
	workers.Watcher = srv.MakeWatcher(workers, DEFAULT_WATCHER_POLLING_INTERVAL)
	if migrationsInterface := srv.Migrations; migrationsInterface != nil {
		workers.Migrations = migrationsInterface.MakeWorker()
	}
	if finhelpInterface := srv.FinHelp; finhelpInterface != nil {
		workers.FinHelp = finhelpInterface.MakeWorker()
	}
	return workers
}

func (workers *Workers) Start() *Workers {
	l4g.Info("Запуск работников")

	workers.startOnce.Do(func() {
		if workers.Migrations != nil {
			go workers.Migrations.Run()
		}
		if workers.FinHelp != nil {
			go workers.FinHelp.Run()
		}
		go workers.Watcher.Start()
	})

	workers.listenerId = workers.ConfigService.AddConfigListener(workers.handleConfigChange)

	return workers
}

func (workers *Workers) handleConfigChange(oldConfig *model.Config, newConfig *model.Config) {

}

func (workers *Workers) Stop() *Workers {
	workers.ConfigService.RemoveConfigListener(workers.listenerId)

	workers.Watcher.Stop()
	if workers.Migrations != nil {
		workers.Migrations.Stop()
	}
	if workers.FinHelp != nil {
		workers.FinHelp.Stop()
	}
	l4g.Info("Остановка работников")

	return workers
}
