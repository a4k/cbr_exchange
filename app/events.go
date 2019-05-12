package app

import (
	"fmt"
	"../utils"
	"os/exec"
)

// Получение журнала событий
func (a *App) GetEventsWin() ([]byte, error) {
	cfg := a.Config().LogSettings
	source := cfg.EventServiceName
	cmd := fmt.Sprintf("Get-Eventlog -LogName %s -Source %s | ConvertTo-JSON", source, source)
	out, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Создание журнала событий
func (a *App) CreateEventLog(serviceName string) (error) {
	resource := utils.FindUtilsFile("EventLogMessages.dll")
	cmd := fmt.Sprintf("New-EventLog -Source '%s' -LogName '%s' -MessageResourceFile '%s'", serviceName, serviceName, resource)
	_, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// Добавление в журнала событий
func (a *App) WriteEvent(serviceName string, etype string, eid uint32, msg string) (error) {
	cmd := fmt.Sprintf("Write-EventLog -LogName '%s' -Source '%s' -EventID %d -EntryType %s -Message '%s'", serviceName, serviceName, eid, etype, msg)
	_, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// Очистка журнала событий
func (a *App) ClearEvents() (error) {
	cfg := a.Config().LogSettings
	source := cfg.EventServiceName
	cmd := fmt.Sprintf("Clear-EventLog -LogName %s", source)
	_, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// Удаление журнала
func RemoveEventLog(logname string) ([]byte, error) {
	cmd := fmt.Sprintf("Remove-EventLog -LogName '%s' -MessageResourceFile '%s'", logname)
	out, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}