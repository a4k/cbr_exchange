package app

import (
	"context"
	"net/http"
	"time"

	l4g "../utils/log4go"
	"../store"
)

type Server struct {
	Store store.Store
	Server      *http.Server

	didFinishListen chan struct{}
}

type RecoveryLogger struct {
}

func (rl *RecoveryLogger) Println(i ...interface{}) {
	l4g.Error("Please check the std error output for the stack trace")
	l4g.Error(i)
}

const TIME_TO_WAIT_FOR_CONNECTIONS_TO_CLOSE_ON_SERVER_SHUTDOWN = time.Second

func (a *App) StartServer() error {
	l4g.Info(("api.server.start_server.starting.info"))

	a.Srv.didFinishListen = make(chan struct{})
	go func() {
		var err error
		if err != nil && err != http.ErrServerClosed {
			l4g.Critical(("api.server.start_server.starting.critical"), err)
			time.Sleep(time.Second)
		}
		close(a.Srv.didFinishListen)
	}()

	return nil
}

func (a *App) StopServer() {
	if a.Srv.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), TIME_TO_WAIT_FOR_CONNECTIONS_TO_CLOSE_ON_SERVER_SHUTDOWN)
		defer cancel()
		didShutdown := false
		for a.Srv.didFinishListen != nil && !didShutdown {
			if err := a.Srv.Server.Shutdown(ctx); err != nil {
				l4g.Warn(err.Error())
			}
			timer := time.NewTimer(time.Millisecond * 50)
			select {
			case <-a.Srv.didFinishListen:
				didShutdown = true
			case <-timer.C:
			}
			timer.Stop()
		}
		a.Srv.Server.Close()
		a.Srv.Server = nil
	}
}
