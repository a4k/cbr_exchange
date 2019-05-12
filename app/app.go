package app

import (
	"crypto/ecdsa"
	"sync/atomic"

	l4g "../utils/log4go"
	"../jobs"
	"../model"

	"../store"
	"../store/sqlstore"
	"../utils"
	tjobs "../jobs/interfaces"
	"../einterfaces"
)

type App struct {
	goroutineCount      int32
	goroutineExitSignal chan struct{}

	Srv *Server
	Jobs    *jobs.JobServer
	FinHelp einterfaces.FinHelpInterface
	newStore func() store.Store

	HubsStopCheckingForDeadlock chan bool

	config          atomic.Value
	configFile      string
	configListeners map[string]func(*model.Config, *model.Config)

	licenseValue       atomic.Value
	clientLicenseValue atomic.Value
	licenseListeners   map[string]func()

	timezones atomic.Value

	configListenerId     string
	licenseListenerId    string
	disableConfigWatch   bool
	configWatcher        *utils.ConfigWatcher
	asymmetricSigningKey *ecdsa.PrivateKey

	clientConfig     map[string]string
	clientConfigHash string
	diagnosticId     string
}

var appCount = 0

// New creates a new App. You must call Shutdown when you're done with it.
// XXX: For now, only one at a time is allowed as some resources are still shared.
func New(options ...Option) (outApp *App, outErr error) {
	appCount++
	if appCount > 1 {
		panic("Only one App should exist at a time. Did you forget to call Shutdown()?")
	}

	app := &App{
		goroutineExitSignal: make(chan struct{}, 1),
		Srv: &Server{
		},
		configFile:       "config.json",
		configListeners:  make(map[string]func(*model.Config, *model.Config)),
		clientConfig:     make(map[string]string),
		licenseListeners: map[string]func(){},
	}
	defer func() {
		if outErr != nil {
			app.Shutdown()
		}
	}()

	for _, option := range options {
		option(app)
	}

	if err := app.LoadConfig(app.configFile); err != nil {
		return nil, err
	}

	app.EnableConfigWatch()

	app.configListenerId = app.AddConfigListener(func(_, _ *model.Config) {
		app.configOrLicenseListener()
	})

	app.regenerateClientConfig()

	l4g.Info(("api.server.new_server.init.info"))

	if app.newStore == nil {
		app.newStore = func() store.Store {
			return store.NewLayeredStore(sqlstore.NewSqlSupplier(app.Config().SqlSettings))
		}
	}

	app.Srv.Store = app.newStore()

	// Создание журнала событий
	err := app.CreateEventLog(app.Config().LogSettings.EventServiceName)
	if err != nil {
		l4g.Info("Error in create event log :: %s", err.Error())
	}

	app.initJobs()

	return app, nil
}

func (a *App) configOrLicenseListener() {
	a.regenerateClientConfig()

}

func (a *App) Shutdown() {
	appCount--

	l4g.Info(("api.server.stop_server.stopping.info"))

	a.StopServer()

	a.WaitForGoroutines()

	if a.Srv.Store != nil {
		a.Srv.Store.Close()
	}
	a.Srv = nil

	a.RemoveConfigListener(a.configListenerId)

	l4g.Info(("api.server.stop_server.stopped.info"))

	a.DisableConfigWatch()
}

func (a *App) initJobs() {

	a.Jobs = jobs.NewJobServer(a, a.Srv.Store)

	if jobsMigrationsInterface != nil {
		a.Jobs.Migrations = jobsMigrationsInterface(a)
	}
	if jobsFinHelpInterface != nil {
		a.Jobs.FinHelp = jobsFinHelpInterface(a)
	}
	a.Jobs.Workers = a.Jobs.InitWorkers()
	a.Jobs.Schedulers = a.Jobs.InitSchedulers()

}

func (a *App) DiagnosticId() string {
	return a.diagnosticId
}

func (a *App) SetDiagnosticId(id string) {
	a.diagnosticId = id
}

func (a *App) EnsureDiagnosticId() {
	if a.diagnosticId != "" {
		return
	}
	a.diagnosticId = ""
}

// Go creates a goroutine, but maintains a record of it to ensure that execution completes before
// the app is destroyed.
func (a *App) Go(f func()) {
	atomic.AddInt32(&a.goroutineCount, 1)

	go func() {
		f()

		atomic.AddInt32(&a.goroutineCount, -1)
		select {
		case a.goroutineExitSignal <- struct{}{}:
		default:
		}
	}()
}

// WaitForGoroutines blocks until all goroutines created by App.Go exit.
func (a *App) WaitForGoroutines() {
	for atomic.LoadInt32(&a.goroutineCount) != 0 {
		<-a.goroutineExitSignal
	}
}


var jobsMigrationsInterface func(*App) tjobs.MigrationsJobInterface

func RegisterJobsMigrationsJobInterface(f func(*App) tjobs.MigrationsJobInterface) {
	jobsMigrationsInterface = f
}

var jobsFinHelpInterface func(*App) tjobs.FinHelpJobInterface

func RegisterJobsFinHelpJobInterface(f func(*App) tjobs.FinHelpJobInterface) {
	jobsFinHelpInterface = f
}

var finhelpInterface func(*App) einterfaces.FinHelpInterface

func RegisterFinHelpInterface(f func(*App) einterfaces.FinHelpInterface) {
	finhelpInterface = f
}
