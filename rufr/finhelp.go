package rufr

import (
	"../app"
	tjobs "../jobs/interfaces"
	"../model"
	"../store"
)

const (
	FINHELP_STATE_UNSCHEDULED = "unscheduled"
	FINHELP_STATE_IN_PROGRESS = "in_progress"
	FINHELP_STATE_COMPLETED   = "completed"

	JOB_DATA_KEY_FINHELP           = "finhelp_key"
	JOB_DATA_KEY_FINHELP_LAST_DONE = "last_done"
)

type FinHelpJobInterfaceImpl struct {
	App *app.App
}

func init() {
	app.RegisterJobsFinHelpJobInterface(func(a *app.App) tjobs.FinHelpJobInterface {
		return &FinHelpJobInterfaceImpl{a}
	})
}

func MakeFinHelpList() []string {
	return []string{
		model.FINHELP_KEY_ADVANCED_PARSING,
	}
}

func GetFinHelpState(finhelp string, store store.Store) (string, *model.Job, *model.AppError) {

	if result := <-store.Job().GetAllByType(model.JOB_TYPE_FINHELP); result.Err != nil {
		return "", nil, result.Err
	} else {
		for _, job := range result.Data.([]*model.Job) {
			if key, ok := job.Data[JOB_DATA_KEY_FINHELP]; ok {
				if key != finhelp {
					continue
				}

				switch job.Status {
				case model.JOB_STATUS_IN_PROGRESS, model.JOB_STATUS_PENDING:
					return FINHELP_STATE_IN_PROGRESS, job, nil
				default:
					return FINHELP_STATE_UNSCHEDULED, job, nil
				}
			}
		}
	}

	return FINHELP_STATE_UNSCHEDULED, nil, nil
}
