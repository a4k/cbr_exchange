package einterfaces

import (
	"context"

	"../model"
)

type FinHelpInterface interface {
	StartSynchronizeJob(ctx context.Context, exportFromTimestamp int64) (*model.Job, *model.AppError)
	managerTestFinHelp()
}
