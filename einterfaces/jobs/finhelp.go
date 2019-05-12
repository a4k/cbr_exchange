package jobs

import (
	"../../model"
)

type FinHelpJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
