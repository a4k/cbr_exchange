package jobs

import (
	"../model"
	"../store"
	tjobs "../jobs/interfaces"
)

type ConfigService interface {
	Config() *model.Config
	AddConfigListener(func(old, current *model.Config)) string
	RemoveConfigListener(string)
}

type StaticConfigService struct {
	Cfg *model.Config
}

func (s StaticConfigService) Config() *model.Config                                   { return s.Cfg }
func (StaticConfigService) AddConfigListener(func(old, current *model.Config)) string { return "" }
func (StaticConfigService) RemoveConfigListener(string)                               {}

type JobServer struct {
	ConfigService ConfigService
	Store         store.Store
	Workers       *Workers
	Schedulers    *Schedulers
	FinHelp       tjobs.FinHelpJobInterface

	Migrations tjobs.MigrationsJobInterface
}

func NewJobServer(configService ConfigService, store store.Store) *JobServer {
	return &JobServer{
		ConfigService: configService,
		Store:         store,
	}
}

func (srv *JobServer) Config() *model.Config {
	return srv.ConfigService.Config()
}

func (srv *JobServer) StartWorkers() {
	srv.Workers = srv.InitWorkers().Start()
}

func (srv *JobServer) StartSchedulers() {
	srv.Schedulers = srv.InitSchedulers().Start()
}

func (srv *JobServer) StopWorkers() {
	if srv.Workers != nil {
		srv.Workers.Stop()
	}
}

func (srv *JobServer) StopSchedulers() {
	if srv.Schedulers != nil {
		srv.Schedulers.Stop()
	}
}
