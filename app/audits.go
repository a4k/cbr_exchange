package app

import (
	l4g "../utils/log4go"
	"fmt"
)
// Вывод в консоль. Логирование RUFR
func (a *App) CLog(eid uint32, text string, v ...interface{}) {
	cfg := a.Config().LogSettings
	source := cfg.EventServiceName
	text2 := fmt.Sprintf("[Код: %d] ", eid)
	text = text2 + text
	l4g.Info(fmt.Sprintf(text, v...))
	err := a.WriteEvent(source, "Information", eid, fmt.Sprintf(text, v...))
	if err != nil {
		l4g.Info("При добавлении события в журнал возникли ошибки %s", err.Error())
	}
}

// Вывод в консоль. Логирование RUFR
func (a *App) CLogErr(eid uint32, text string, v ...interface{}) {
	cfg := a.Config().LogSettings
	source := cfg.EventServiceName
	text2 := fmt.Sprintf("[Код: %d] ", eid)
	text = text2 + text
	l4g.Error(fmt.Sprintf(text, v...))
	err := a.WriteEvent(source, "Error", eid, fmt.Sprintf(text, v...))
	if err != nil {
		l4g.Info("При добавлении события в журнал возникли ошибки %s", err.Error())
	}
}

// Вывод в консоль. Логирование RUFR
func (a *App) CLogWarning(eid uint32, text string, v ...interface{}) {
	cfg := a.Config().LogSettings
	source := cfg.EventServiceName
	text2 := fmt.Sprintf("[Код: %d] ", eid)
	text = text2 + text
	l4g.Warn(fmt.Sprintf(text, v...))
	err := a.WriteEvent(source, "Warning", eid, fmt.Sprintf(text, v...))
	if err != nil {
		l4g.Info("При добавлении события в журнал возникли ошибки %s", err.Error())
	}
}