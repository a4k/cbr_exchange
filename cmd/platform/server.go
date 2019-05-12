package main

import (
	"os"

	"github.com/judwhite/go-svc/svc"
	l4g "../../utils/log4go"

	_ "../../rufr"
	"../../app"

	"../../utils"

	"os/signal"
	"syscall"

	"sync"
)

type server struct {
	data chan int

	exit chan struct{}
	wg   sync.WaitGroup
}

// implements svc.Service
type program struct {
	svr *server
}

func main() {
	prg := program{
		svr: &server{},
	}

	// call svc.Run to start your program/service
	// svc.Run will call Init, Start, and Stop
	if err := svc.Run(&prg); err != nil {
		l4g.Crash(err)
	}
}


func (p *program) Init(env svc.Environment) error {
	l4g.Debug("is win service? %v\n", env.IsWindowsService())

	if env.IsWindowsService() {


	}

	return nil
}

func (p *program) Start() error {
	l4g.Debug("Starting...\n")
	go p.svr.start()
	return nil
}

func (p *program) Stop() error {
	l4g.Debug("Stopping...\n")
	if err := p.svr.stop(); err != nil {
		return err
	}
	l4g.Debug("Stopped.\n")
	return nil
}

func (s *server) start() {
	s.data = make(chan int)
	s.exit = make(chan struct{})

	s.wg.Add(2)

	go func() {
		interruptChan := make(chan os.Signal, 1)
		runServer("default.json", false, interruptChan)

	}()

}

func (s *server) stop() error {
	close(s.exit)
	//s.wg.Wait()
	return nil
}

func runServer(configFileLocation string, disableConfigWatch bool, interruptChan chan os.Signal) error {
	options := []app.Option{app.ConfigFile(configFileLocation)}
	if disableConfigWatch {
		options = append(options, app.DisableConfigWatch)
	}

	a, err := app.New(options...)
	if err != nil {
		l4g.Critical(err.Error())
		return err
	}
	defer a.Shutdown()

	pwd, _ := os.Getwd()

	l4g.Info(("working_dir %s "), pwd)
	l4g.Info(("config_file %s"), utils.FindConfigFile(configFileLocation))

	serverErr := a.StartServer()
	if serverErr != nil {
		l4g.Critical(serverErr.Error())
		return serverErr
	}

	a.ReloadConfig()

	if *a.Config().JobSettings.RunJobs {
		a.Jobs.StartWorkers()
	}
	if *a.Config().JobSettings.RunScheduler {
		a.Jobs.StartSchedulers()
	}

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	a.Jobs.StopSchedulers()
	a.Jobs.StopWorkers()

	return nil
}

