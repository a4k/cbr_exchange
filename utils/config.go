package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	l4g "../utils/log4go"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"net/http"

	"../model"
)

const (
	LOG_ROTATE_SIZE = 10000
	LOG_FILENAME    = "server.log"
)

var originalDisableDebugLvl l4g.Level = l4g.DEBUG

// FindConfigFile attempts to find an existing configuration file. fileName can be an absolute or
// relative path or name such as "/opt/mattermost/config.json" or simply "config.json". An empty
// string is returned if no configuration is found.
func FindConfigFile(fileName string) (path string) {
	if filepath.IsAbs(fileName) {
		if _, err := os.Stat(fileName); err == nil {
			return fileName
		}
	} else {
		for _, dir := range []string{"./config", "../config", "../../config", "."} {
			path, _ := filepath.Abs(filepath.Join(dir, fileName))
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}
	return ""
}
func FindUtilsFile(fileName string) (path string) {
	if filepath.IsAbs(fileName) {
		if _, err := os.Stat(fileName); err == nil {
			return fileName
		}
	} else {
		for _, dir := range []string{"./utils", "../utils", "../../utils", "."} {
			path, _ := filepath.Abs(filepath.Join(dir, fileName))
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}
	return ""
}
func FindExeFile(fileName string) (path string) {
	if filepath.IsAbs(fileName) {
		if _, err := os.Stat(fileName); err == nil {
			return fileName
		}
	} else {
		for _, dir := range []string{"./utils/arj", "../utils/arj", "../../utils/arj", "."} {
			path, _ := filepath.Abs(filepath.Join(dir, fileName))
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}
	return ""
}

// FindDir looks for the given directory in nearby ancestors, falling back to `./` if not found.
func FindDir(dir string) (string, bool) {
	for _, parent := range []string{".", "..", "../.."} {
		foundDir, err := filepath.Abs(filepath.Join(parent, dir))
		if err != nil {
			continue
		} else if _, err := os.Stat(foundDir); err == nil {
			return foundDir, true
		}
	}
	return "./", false
}

func DisableDebugLogForTest() {
	if l4g.Global["stdout"] != nil {
		originalDisableDebugLvl = l4g.Global["stdout"].Level
		l4g.Global["stdout"].Level = l4g.ERROR
	}
}

func EnableDebugLogForTest() {
	if l4g.Global["stdout"] != nil {
		l4g.Global["stdout"].Level = originalDisableDebugLvl
	}
}

func ConfigureCmdLineLog() {
	ls := model.LogSettings{}
	ls.EnableConsole = true
	ls.ConsoleLevel = "WARN"
	ConfigureLog(&ls)
}

// TODO: this code initializes console and file logging. It will eventually be replaced by JSON logging in logger/logger.go
// See PLT-3893 for more information
func ConfigureLog(s *model.LogSettings) {

	l4g.Close()

	if s.EnableConsole {
		level := l4g.DEBUG
		if s.ConsoleLevel == "INFO" {
			level = l4g.INFO
		} else if s.ConsoleLevel == "WARN" {
			level = l4g.WARNING
		} else if s.ConsoleLevel == "ERROR" {
			level = l4g.ERROR
		}

		lw := l4g.NewConsoleLogWriter()
		lw.SetFormat("[%D %T] [%L] %M")
		l4g.AddFilter("stdout", level, lw)
	}

	if s.EnableFile {

		var fileFormat = s.FileFormat

		if fileFormat == "" {
			fileFormat = "[%D %T] [%L] %M"
		}

		level := l4g.DEBUG
		if s.FileLevel == "INFO" {
			level = l4g.INFO
		} else if s.FileLevel == "WARN" {
			level = l4g.WARNING
		} else if s.FileLevel == "ERROR" {
			level = l4g.ERROR
		}

		flw := l4g.NewFileLogWriter(GetLogFileLocation(s.FileLocation), false, s.EventServiceName)
		flw.SetFormat(fileFormat)
		flw.SetRotate(true)
		flw.SetRotateLines(LOG_ROTATE_SIZE)
		l4g.AddFilter("file", level, flw)
	}
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = FindDir("logs")
	}

	return filepath.Join(fileLocation, LOG_FILENAME)
}

func SaveConfig(fileName string, config *model.Config) *model.AppError {
	b, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return model.NewAppError("SaveConfig", "utils.config.save_config.saving.app_error",
			map[string]interface{}{"Filename": fileName}, err.Error(), http.StatusBadRequest)
	}

	err = ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		return model.NewAppError("SaveConfig", "utils.config.save_config.saving.app_error",
			map[string]interface{}{"Filename": fileName}, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

type ConfigWatcher struct {
	watcher *fsnotify.Watcher
	close   chan struct{}
	closed  chan struct{}
}

func NewConfigWatcher(cfgFileName string, f func()) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create config watcher for file: "+cfgFileName)
	}

	configFile := filepath.Clean(cfgFileName)
	configDir, _ := filepath.Split(configFile)
	watcher.Add(configDir)

	ret := &ConfigWatcher{
		watcher: watcher,
		close:   make(chan struct{}),
		closed:  make(chan struct{}),
	}

	go func() {
		defer close(ret.closed)
		defer watcher.Close()

		for {
			select {
			case event := <-watcher.Events:
				// we only care about the config file
				if filepath.Clean(event.Name) == configFile {
					if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
						l4g.Info(fmt.Sprintf("Config file watcher detected a change reloading %v", cfgFileName))

						if _, configReadErr := ReadConfigFile(cfgFileName, true); configReadErr == nil {
							f()
						} else {
							l4g.Error(fmt.Sprintf("Failed to read while watching config file at %v with err=%v", cfgFileName, configReadErr.Error()))
						}
					}
				}
			case err := <-watcher.Errors:
				l4g.Error(fmt.Sprintf("Failed while watching config file at %v with err=%v", cfgFileName, err.Error()))
			case <-ret.close:
				return
			}
		}
	}()

	return ret, nil
}

func (w *ConfigWatcher) Close() {
	close(w.close)
	<-w.closed
}

// ReadConfig reads and parses the given configuration.
func ReadConfig(r io.Reader, allowEnvironmentOverrides bool) (*model.Config, error) {
	v := viper.New()

	if allowEnvironmentOverrides {
		v.SetEnvPrefix("mm")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()
	}

	v.SetConfigType("json")
	if err := v.ReadConfig(r); err != nil {
		return nil, err
	}

	var config model.Config
	unmarshalErr := v.Unmarshal(&config)
	if unmarshalErr == nil {
		// https://github.com/spf13/viper/issues/324
		// https://github.com/spf13/viper/issues/348

	}
	return &config, unmarshalErr
}

// ReadConfigFile reads and parses the configuration at the given file path.
func ReadConfigFile(path string, allowEnvironmentOverrides bool) (*model.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadConfig(f, allowEnvironmentOverrides)
}

// EnsureConfigFile will attempt to locate a config file with the given name. If it does not exist,
// it will attempt to locate a default config file, and copy it to a file named fileName in the same
// directory. In either case, the config file path is returned.
func EnsureConfigFile(fileName string) (string, error) {
	if configFile := FindConfigFile(fileName); configFile != "" {
		return configFile, nil
	}
	if defaultPath := FindConfigFile("default.json"); defaultPath != "" {
		destPath := filepath.Join(filepath.Dir(defaultPath), fileName)
		src, err := os.Open(defaultPath)
		if err != nil {
			return "", err
		}
		defer src.Close()
		dest, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return "", err
		}
		defer dest.Close()
		if _, err := io.Copy(dest, src); err == nil {
			return destPath, nil
		}
	}
	return "", fmt.Errorf("no config file found")
}

// LoadConfig will try to search around for the corresponding config file.  It will search
// /tmp/fileName then attempt ./config/fileName, then ../config/fileName and last it will look at
// fileName.
func LoadConfig(fileName string) (config *model.Config, configPath string, appErr *model.AppError) {
	if fileName != filepath.Base(fileName) {
		configPath = fileName
	} else {
		if path, err := EnsureConfigFile(fileName); err != nil {
			appErr = model.NewAppError("LoadConfig", "utils.config.load_config.opening.panic", map[string]interface{}{"Filename": fileName, "Error": err.Error()}, "", 0)
			return
		} else {
			configPath = path
		}
	}

	config, err := ReadConfigFile(configPath, true)
	if err != nil {
		appErr = model.NewAppError("LoadConfig", "utils.config.load_config.decoding.panic", map[string]interface{}{"Filename": fileName, "Error": err.Error()}, "", 0)
		return
	}

	config.SetDefaults()

	if err := config.IsValid(); err != nil {
		return nil, "", err
	}

	return config, configPath, nil
}

func GenerateClientConfig(c *model.Config, diagnosticId string) map[string]string {
	props := make(map[string]string)

	props["SQLDriverName"] = *c.SqlSettings.DriverName

	props["DiagnosticId"] = diagnosticId
	props["DiagnosticsEnabled"] = strconv.FormatBool(*c.LogSettings.EnableDiagnostics)

	// Set default values for all options that require a license.
	props["ExperimentalTownSquareIsReadOnly"] = "false"
	props["ExperimentalEnableAuthenticationTransfer"] = "true"
	props["EnableCustomBrand"] = "false"
	props["CustomBrandText"] = ""
	props["CustomDescriptionText"] = ""
	props["EnableLdap"] = "false"
	props["LdapLoginFieldName"] = ""
	props["LdapNicknameAttributeSet"] = "false"
	props["LdapFirstNameAttributeSet"] = "false"
	props["LdapLastNameAttributeSet"] = "false"
	props["LdapLoginButtonColor"] = ""
	props["LdapLoginButtonBorderColor"] = ""
	props["LdapLoginButtonTextColor"] = ""
	props["EnableMultifactorAuthentication"] = "false"
	props["EnforceMultifactorAuthentication"] = "false"
	props["EnableCompliance"] = "false"
	props["EnableMobileFileDownload"] = "true"
	props["EnableMobileFileUpload"] = "true"
	props["EnableSaml"] = "false"
	props["SamlLoginButtonText"] = ""
	props["SamlFirstNameAttributeSet"] = "false"
	props["SamlLastNameAttributeSet"] = "false"
	props["SamlNicknameAttributeSet"] = "false"
	props["SamlLoginButtonColor"] = ""
	props["SamlLoginButtonBorderColor"] = ""
	props["SamlLoginButtonTextColor"] = ""
	props["EnableCluster"] = "false"
	props["EnableMetrics"] = "false"
	props["EnableSignUpWithGoogle"] = "false"
	props["EnableSignUpWithOffice365"] = "false"
	props["PasswordMinimumLength"] = "0"
	props["PasswordRequireLowercase"] = "false"
	props["PasswordRequireUppercase"] = "false"
	props["PasswordRequireNumber"] = "false"
	props["PasswordRequireSymbol"] = "false"
	props["EnableBanner"] = "false"
	props["BannerText"] = ""
	props["BannerColor"] = ""
	props["BannerTextColor"] = ""
	props["AllowBannerDismissal"] = "false"
	props["EnableThemeSelection"] = "true"
	props["DefaultTheme"] = ""
	props["AllowCustomThemes"] = "true"
	props["AllowedThemes"] = ""
	props["DataRetentionEnableMessageDeletion"] = "false"
	props["DataRetentionMessageRetentionDays"] = "0"
	props["DataRetentionEnableFileDeletion"] = "false"
	props["DataRetentionFileRetentionDays"] = "0"

	return props
}

