package jobs

import (
	"math/rand"
	"time"

	l4g "../utils/log4go"
	"../model"
	"fmt"
)

// Default polling interval for jobs termination.
// (Defining as `var` rather than `const` allows tests to lower the interval.)
var DEFAULT_WATCHER_POLLING_INTERVAL = 15000

type Watcher struct {
	srv     *JobServer
	workers *Workers

	stop            chan bool
	stopped         chan bool
	pollingInterval int
}

func (srv *JobServer) MakeWatcher(workers *Workers, pollingInterval int) *Watcher {
	return &Watcher{
		stop:            make(chan bool, 1),
		stopped:         make(chan bool, 1),
		pollingInterval: pollingInterval,
		workers:         workers,
		srv:             srv,
	}
}

func (watcher *Watcher) Start() {
	l4g.Debug("Наблюдатель запустился")

	// Delay for some random number of milliseconds before starting to ensure that multiple
	// instances of the jobserver  don't poll at a time too close to each other.
	rand.Seed(time.Now().UTC().UnixNano())
	<-time.After(time.Duration(rand.Intn(watcher.pollingInterval)) * time.Millisecond)

	defer func() {
		l4g.Debug("Watcher Finished")
		watcher.stopped <- true
	}()

	for {
		select {
		case <-watcher.stop:
			l4g.Debug("Watcher: Received stop signal")
			return
		case <-time.After(time.Duration(watcher.pollingInterval) * time.Millisecond):
			watcher.PollAndNotify()
		}
	}
}

func (watcher *Watcher) Stop() {
	l4g.Debug("Watcher Stopping")
	watcher.stop <- true
	<-watcher.stopped
}

func (watcher *Watcher) PollAndNotify() {
	if result := <-watcher.srv.Store.Job().GetAllByStatus(model.JOB_STATUS_PENDING); result.Err != nil {
		l4g.Error(fmt.Sprintf("Error occurred getting all pending statuses: %v", result.Err.Error()))
	} else {
		jobs := result.Data.([]*model.Job)

		for _, job := range jobs {
			if job.Type == model.JOB_TYPE_MIGRATIONS {
				if watcher.workers.Migrations != nil {
					select {
					case watcher.workers.Migrations.JobChannel() <- *job:
					default:
					}
				}
			}
			if job.Type == model.JOB_TYPE_FINHELP {
				if watcher.workers.FinHelp != nil {
					select {
					case watcher.workers.FinHelp.JobChannel() <- *job:
					default:
					}
				}
			}
		}
	}
}
