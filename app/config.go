package app

import (
	"crypto/md5"

	"encoding/json"
	"fmt"
	"net/url"
	"runtime/debug"

	l4g "../utils/log4go"

	"../model"
	"../utils"
)

func (a *App) Config() *model.Config {
	if cfg := a.config.Load(); cfg != nil {
		return cfg.(*model.Config)
	}
	return &model.Config{}
}

func (a *App) UpdateConfig(f func(*model.Config)) {
	old := a.Config()
	updated := old.Clone()
	f(updated)
	a.config.Store(updated)

	a.InvokeConfigListeners(old, updated)
}

func (a *App) PersistConfig() {
	utils.SaveConfig(a.ConfigFileName(), a.Config())
}

func (a *App) LoadConfig(configFile string) *model.AppError {
	old := a.Config()

	cfg, configPath, err := utils.LoadConfig(configFile)
	if err != nil {
		return err
	}

	a.configFile = configPath

	utils.ConfigureLog(&cfg.LogSettings)

	a.config.Store(cfg)

	a.InvokeConfigListeners(old, cfg)
	return nil
}

func (a *App) ReloadConfig() *model.AppError {
	debug.FreeOSMemory()
	if err := a.LoadConfig(a.configFile); err != nil {
		return err
	}

	return nil
}

func (a *App) ConfigFileName() string {
	return a.configFile
}

func (a *App) ClientConfig() map[string]string {
	return a.clientConfig
}

func (a *App) ClientConfigHash() string {
	return a.clientConfigHash
}

func (a *App) EnableConfigWatch() {
	if a.configWatcher == nil && !a.disableConfigWatch {
		configWatcher, err := utils.NewConfigWatcher(a.ConfigFileName(), func() {
			a.ReloadConfig()
		})
		if err != nil {
			l4g.Error(err)
		}
		a.configWatcher = configWatcher
	}
}

func (a *App) DisableConfigWatch() {
	if a.configWatcher != nil {
		a.configWatcher.Close()
		a.configWatcher = nil
	}
}

// Registers a function with a given to be called when the config is reloaded and may have changed. The function
// will be called with two arguments: the old config and the new config. AddConfigListener returns a unique ID
// for the listener that can later be used to remove it.
func (a *App) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	id := model.NewId()
	a.configListeners[id] = listener
	return id
}

// Removes a listener function by the unique ID returned when AddConfigListener was called
func (a *App) RemoveConfigListener(id string) {
	delete(a.configListeners, id)
}

func (a *App) InvokeConfigListeners(old, current *model.Config) {
	for _, listener := range a.configListeners {
		listener(old, current)
	}
}

func (a *App) regenerateClientConfig() {
	a.clientConfig = utils.GenerateClientConfig(a.Config(), a.DiagnosticId())

	clientConfigJSON, _ := json.Marshal(a.clientConfig)
	a.clientConfigHash = fmt.Sprintf("%x", md5.Sum(clientConfigJSON))
}

func (a *App) Desanitize(cfg *model.Config) {
	actual := a.Config()

	if *cfg.SqlSettings.DataSource == model.FAKE_SETTING {
		*cfg.SqlSettings.DataSource = *actual.SqlSettings.DataSource
	}
	if cfg.SqlSettings.AtRestEncryptKey == model.FAKE_SETTING {
		cfg.SqlSettings.AtRestEncryptKey = actual.SqlSettings.AtRestEncryptKey
	}

	for i := range cfg.SqlSettings.DataSourceReplicas {
		cfg.SqlSettings.DataSourceReplicas[i] = actual.SqlSettings.DataSourceReplicas[i]
	}

	for i := range cfg.SqlSettings.DataSourceSearchReplicas {
		cfg.SqlSettings.DataSourceSearchReplicas[i] = actual.SqlSettings.DataSourceSearchReplicas[i]
	}
}

func (a *App) GetCookieDomain() string {
	if *a.Config().ServiceSettings.AllowCookiesForSubdomains {
		if siteURL, err := url.Parse(*a.Config().ServiceSettings.SiteURL); err == nil {
			return siteURL.Hostname()
		}
	}
	return ""
}

// ClientConfigWithNoAccounts gets the configuration in a format suitable for sending to the client.
func (a *App) ClientConfigWithNoAccounts() map[string]string {
	respCfg := map[string]string{}
	for k, v := range a.ClientConfig() {
		respCfg[k] = v
	}

	return respCfg
}
