package model

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	DATABASE_DRIVER_MYSQL    = "mysql"
	DATABASE_DRIVER_MSSQL    = "mssql"
	DATABASE_DRIVER_POSTGRES = "postgres"

	GENERIC_NO_CHANNEL_NOTIFICATION = "generic_no_channel"
	GENERIC_NOTIFICATION            = "generic"
	FULL_NOTIFICATION               = "full"

	FAKE_SETTING = "********************************"

	SQL_SETTINGS_DEFAULT_DATA_SOURCE = "mmuser:mostest@tcp(dockerhost:3306)/dev?charset=utf8mb4,utf8&readTimeout=30s&writeTimeout=30s"

	CONN_SECURITY_NONE     = ""
	CONN_SECURITY_PLAIN    = "PLAIN"
	CONN_SECURITY_TLS      = "TLS"
	CONN_SECURITY_STARTTLS = "STARTTLS"

	EMAIL_BATCHING_BUFFER_SIZE = 256
	EMAIL_BATCHING_INTERVAL    = 30

	EMAIL_NOTIFICATION_CONTENTS_FULL             = "full"
	EMAIL_NOTIFICATION_CONTENTS_GENERIC          = "generic"
	EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION = ""

	SERVICE_SETTINGS_DEFAULT_SITE_URL           = ""
	SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE      = ""
	SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE       = ""
	SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT       = 300
	SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT      = 300
	SERVICE_SETTINGS_DEFAULT_MAX_LOGIN_ATTEMPTS = 10
	SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM    = ""
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":8065"

	SHOW_USERNAME          = "username"
	SHOW_NICKNAME_FULLNAME = "nickname_full_name"
	SHOW_FULLNAME          = "full_name"

	FINHELP_SETTINGS_DEFAULT_REQUEST_TIMEOUT_SECONDS = 60
	FINHELP_SETTINGS_DEFAULT_JOB_START_TIME          = "03:00"
	FINHELP_SETTINGS_DEFAULT_SUBJECT				 = "[#rufr]"
	FINHELP_SETTINGS_DEFAULT_DIR_TO_DOCTEMP				 = "demo"
	FINHELP_SETTINGS_DEFAULT_DOCFILE_BANK       = "BANC.DBF"
	FINHELP_SETTINGS_DEFAULT_DOCFILE_MICROFIN   = "microfin.csv"
	FINHELP_SETTINGS_DEFAULT_DOCFILE_B_ARX      = "B_ARX.DBF"
	FINHELP_SETTINGS_DEFAULT_DOCFILE_STRAHOV    = "strahov.csv"
)

type SqlSettings struct {
	DriverName               *string
	DataSource               *string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	MaxIdleConns             *int
	MaxOpenConns             *int
	Trace                    bool
	AtRestEncryptKey         string
	QueryTimeout             *int
}

func (s *SqlSettings) SetDefaults() {
	if s.DriverName == nil {
		s.DriverName = NewString(DATABASE_DRIVER_MYSQL)
	}

	if s.DataSource == nil {
		s.DataSource = NewString(SQL_SETTINGS_DEFAULT_DATA_SOURCE)
	}

	if len(s.AtRestEncryptKey) == 0 {
		s.AtRestEncryptKey = NewRandomString(32)
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewInt(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = NewInt(300)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewInt(30)
	}
}

type LogSettings struct {
	EnableConsole          bool
	ConsoleLevel           string
	EnableFile             bool
	FileLevel              string
	FileFormat             string
	FileLocation           string
	EventServiceName          string
	EnableWebhookDebugging bool
	EnableDiagnostics      *bool
}

func (s *LogSettings) SetDefaults() {
	if s.EnableDiagnostics == nil {
		s.EnableDiagnostics = NewBool(true)
	}
}

type ConfigFunc func() *Config

type Config struct {
	SqlSettings             SqlSettings
	LogSettings             LogSettings
	ServiceSettings         ServiceSettings
	RateLimitSettings       RateLimitSettings
	ImportInventorySettings ImportInventorySettings
	ImportUnitSettings      ImportUnitSettings
	JobSettings             JobSettings
	ExchangeSettings         ExchangeSettings
	EmailSettings           EmailSettings
}

func (o *Config) Clone() *Config {
	var ret Config
	if err := json.Unmarshal([]byte(o.ToJson()), &ret); err != nil {
		panic(err)
	}
	return &ret
}

func (o *Config) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func ConfigFromJson(data io.Reader) *Config {
	decoder := json.NewDecoder(data)
	var o Config
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

func (o *Config) SetDefaults() {

	o.SqlSettings.SetDefaults()
	o.ExchangeSettings.SetDefaults()
	o.EmailSettings.SetDefaults()

	o.LogSettings.SetDefaults()

}

func (o *Config) IsValid() *AppError {

	if err := o.SqlSettings.isValid(); err != nil {
		return err
	}
	if err := o.ExchangeSettings.isValid(); err != nil {
		return err
	}
	if err := o.EmailSettings.isValid(); err != nil {
		return err
	}

	return nil
}

func (fin *ExchangeSettings) isValid() *AppError {
	//if *fin.EnableIndexing {
	//	if len(*ess.ConnectionUrl) == 0 {
	//		return NewAppError("Config.IsValid", "model.config.is_valid.elastic_search.connection_url.app_error", nil, "", http.StatusBadRequest)
	//	}
	//}

	if _, err := time.Parse("15:04", *fin.DailyRunTime); err != nil {
		return NewAppError("Config.IsValid", "model.config.is_valid.finhelp.daily_run_time.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return nil
}

type ExchangeSettings struct {
	RequestTimeoutSeconds *int
	DailyRunTime          *string
	Subject          	  string
	DirToDocTemp          string
	Docfile_bank          string
	Docfile_barx          string
	Docfile_strahov       string
	Docfile_microfin      string
	ForceUpdate           *bool
}

func (s *ExchangeSettings) SetDefaults() {
	if s.RequestTimeoutSeconds == nil {
		s.RequestTimeoutSeconds = NewInt(FINHELP_SETTINGS_DEFAULT_REQUEST_TIMEOUT_SECONDS)
	}
	if s.DailyRunTime == nil {
		s.DailyRunTime = NewString(FINHELP_SETTINGS_DEFAULT_JOB_START_TIME)
	}
	if s.ForceUpdate == nil {
		s.ForceUpdate = NewBool(false)
	}
}

func (es *EmailSettings) isValid() *AppError {
	if !(es.ConnectionSecurity == CONN_SECURITY_NONE || es.ConnectionSecurity == CONN_SECURITY_TLS || es.ConnectionSecurity == CONN_SECURITY_STARTTLS || es.ConnectionSecurity == CONN_SECURITY_PLAIN) {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_security.app_error", nil, "", http.StatusBadRequest)
	}

	if len(es.InviteSalt) < 32 {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_salt.app_error", nil, "", http.StatusBadRequest)
	}

	if *es.EmailBatchingBufferSize <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_batching_buffer_size.app_error", nil, "", http.StatusBadRequest)
	}

	if *es.EmailBatchingInterval < 30 {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_batching_interval.app_error", nil, "", http.StatusBadRequest)
	}

	if !(*es.EmailNotificationContentsType == EMAIL_NOTIFICATION_CONTENTS_FULL || *es.EmailNotificationContentsType == EMAIL_NOTIFICATION_CONTENTS_GENERIC) {
		return NewAppError("Config.IsValid", "model.config.is_valid.email_notification_contents_type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type EmailSettings struct {
	EnableSignUpWithEmail             bool
	EnableSignInWithEmail             *bool
	EnableSignInWithUsername          *bool
	SendEmailNotifications            bool
	UseChannelInEmailNotifications    *bool
	RequireEmailVerification          bool
	FeedbackName                      string
	FeedbackEmail                     string
	FeedbackOrganization              *string
	EnableSMTPAuth                    *bool
	SMTPUsername                      string
	SMTPPassword                      string
	SMTPServer                        string
	SMTPPort                          string
	ConnectionSecurity                string
	InviteSalt                        string
	SendPushNotifications             *bool
	PushNotificationServer            *string
	PushNotificationContents          *string
	EnableEmailBatching               *bool
	EmailBatchingBufferSize           *int
	EmailBatchingInterval             *int
	EnablePreviewModeBanner           *bool
	SkipServerCertificateVerification *bool
	EmailNotificationContentsType     *string
	LoginButtonColor                  *string
	LoginButtonBorderColor            *string
	LoginButtonTextColor              *string
	Pop3Login                         *string
	Pop3Password                      *string
	Pop3Url                           *string
}

func (s *EmailSettings) SetDefaults() {
	if len(s.InviteSalt) == 0 {
		s.InviteSalt = NewRandomString(32)
	}

	if s.EnableSignInWithEmail == nil {
		s.EnableSignInWithEmail = NewBool(s.EnableSignUpWithEmail)
	}

	if s.EnableSignInWithUsername == nil {
		s.EnableSignInWithUsername = NewBool(false)
	}

	if s.UseChannelInEmailNotifications == nil {
		s.UseChannelInEmailNotifications = NewBool(false)
	}

	if s.SendPushNotifications == nil {
		s.SendPushNotifications = NewBool(false)
	}

	if s.PushNotificationServer == nil {
		s.PushNotificationServer = NewString("")
	}

	if s.PushNotificationContents == nil {
		s.PushNotificationContents = NewString(GENERIC_NOTIFICATION)
	}

	if s.FeedbackOrganization == nil {
		s.FeedbackOrganization = NewString(EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION)
	}

	if s.EnableEmailBatching == nil {
		s.EnableEmailBatching = NewBool(false)
	}

	if s.EmailBatchingBufferSize == nil {
		s.EmailBatchingBufferSize = NewInt(EMAIL_BATCHING_BUFFER_SIZE)
	}

	if s.EmailBatchingInterval == nil {
		s.EmailBatchingInterval = NewInt(EMAIL_BATCHING_INTERVAL)
	}

	if s.EnablePreviewModeBanner == nil {
		s.EnablePreviewModeBanner = NewBool(true)
	}

	if s.EnableSMTPAuth == nil {
		s.EnableSMTPAuth = new(bool)
		if s.ConnectionSecurity == CONN_SECURITY_NONE {
			*s.EnableSMTPAuth = false
		} else {
			*s.EnableSMTPAuth = true
		}
	}

	if s.ConnectionSecurity == CONN_SECURITY_PLAIN {
		s.ConnectionSecurity = CONN_SECURITY_NONE
	}

	if s.SkipServerCertificateVerification == nil {
		s.SkipServerCertificateVerification = NewBool(false)
	}

	if s.EmailNotificationContentsType == nil {
		s.EmailNotificationContentsType = NewString(EMAIL_NOTIFICATION_CONTENTS_FULL)
	}

	if s.LoginButtonColor == nil {
		s.LoginButtonColor = NewString("#0000")
	}

	if s.LoginButtonBorderColor == nil {
		s.LoginButtonBorderColor = NewString("#2389D7")
	}

	if s.LoginButtonTextColor == nil {
		s.LoginButtonTextColor = NewString("#2389D7")
	}
}

func (ss *SqlSettings) isValid() *AppError {
	if len(ss.AtRestEncryptKey) < 32 {
		return NewAppError("Config.IsValid", "model.config.is_valid.encrypt_sql.app_error", nil, "", http.StatusBadRequest)
	}

	if !(*ss.DriverName == DATABASE_DRIVER_MYSQL || *ss.DriverName == DATABASE_DRIVER_POSTGRES || *ss.DriverName == DATABASE_DRIVER_MSSQL) {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_driver.app_error", nil, "", http.StatusBadRequest)
	}

	if *ss.MaxIdleConns <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_idle.app_error", nil, "", http.StatusBadRequest)
	}

	if *ss.QueryTimeout <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_query_timeout.app_error", nil, "", http.StatusBadRequest)
	}

	if len(*ss.DataSource) == 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_data_src.app_error", nil, "", http.StatusBadRequest)
	}

	if *ss.MaxOpenConns <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.sql_max_conn.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ServiceSettings struct {
	SiteURL                                           *string
	WebsocketURL                                      *string
	LicenseFileLocation                               *string
	ListenAddress                                     *string
	ConnectionSecurity                                *string
	TLSCertFile                                       *string
	TLSKeyFile                                        *string
	UseLetsEncrypt                                    *bool
	LetsEncryptCertificateCacheFile                   *string
	Forward80To443                                    *bool
	ReadTimeout                                       *int
	WriteTimeout                                      *int
	MaximumLoginAttempts                              *int
	GoroutineHealthThreshold                          *int
	GoogleDeveloperKey                                string
	EnableOAuthServiceProvider                        bool
	EnableIncomingWebhooks                            bool
	EnableOutgoingWebhooks                            bool
	EnableCommands                                    *bool
	EnableOnlyAdminIntegrations                       *bool
	EnablePostUsernameOverride                        bool
	EnablePostIconOverride                            bool
	EnableAPIv3                                       *bool
	EnableLinkPreviews                                *bool
	EnableTesting                                     bool
	EnableDeveloper                                   *bool
	EnableSecurityFixAlert                            *bool
	EnableInsecureOutgoingConnections                 *bool
	AllowedUntrustedInternalConnections               *string
	EnableMultifactorAuthentication                   *bool
	EnforceMultifactorAuthentication                  *bool
	EnableUserAccessTokens                            *bool
	AllowCorsFrom                                     *string
	AllowCookiesForSubdomains                         *bool
	SessionLengthWebInDays                            *int
	SessionLengthMobileInDays                         *int
	SessionLengthSSOInDays                            *int
	SessionCacheInMinutes                             *int
	SessionIdleTimeoutInMinutes                       *int
	WebsocketSecurePort                               *int
	WebsocketPort                                     *int
	WebserverMode                                     *string
	EnableCustomEmoji                                 *bool
	EnableEmojiPicker                                 *bool
	RestrictCustomEmojiCreation                       *string
	RestrictPostDelete                                *string
	AllowEditPost                                     *string
	PostEditTimeLimit                                 *int
	TimeBetweenUserTypingUpdatesMilliseconds          *int64
	EnablePostSearch                                  *bool
	EnableUserTypingMessages                          *bool
	EnableChannelViewedMessages                       *bool
	EnableUserStatuses                                *bool
	ExperimentalEnableAuthenticationTransfer          *bool
	ClusterLogTimeoutMilliseconds                     *int
	CloseUnusedDirectMessages                         *bool
	EnablePreviewFeatures                             *bool
	EnableTutorial                                    *bool
	ExperimentalEnableDefaultChannelLeaveJoinMessages *bool
	ExperimentalGroupUnreadChannels                   *string
	ImageProxyType                                    *string
	ImageProxyURL                                     *string
	ImageProxyOptions                                 *string
}

func (s *ServiceSettings) SetDefaults() {
	if s.SiteURL == nil {
		s.SiteURL = NewString(SERVICE_SETTINGS_DEFAULT_SITE_URL)
	}

	if s.WebsocketURL == nil {
		s.WebsocketURL = NewString("")
	}

	if s.LicenseFileLocation == nil {
		s.LicenseFileLocation = NewString("")
	}

	if s.ListenAddress == nil {
		s.ListenAddress = NewString(SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS)
	}

	if s.EnableAPIv3 == nil {
		s.EnableAPIv3 = NewBool(true)
	}

	if s.EnableLinkPreviews == nil {
		s.EnableLinkPreviews = NewBool(false)
	}

	if s.EnableDeveloper == nil {
		s.EnableDeveloper = NewBool(false)
	}

	if s.EnableSecurityFixAlert == nil {
		s.EnableSecurityFixAlert = NewBool(true)
	}

	if s.EnableInsecureOutgoingConnections == nil {
		s.EnableInsecureOutgoingConnections = NewBool(false)
	}

	if s.AllowedUntrustedInternalConnections == nil {
		s.AllowedUntrustedInternalConnections = NewString("")
	}

	if s.EnableMultifactorAuthentication == nil {
		s.EnableMultifactorAuthentication = NewBool(false)
	}

	if s.EnforceMultifactorAuthentication == nil {
		s.EnforceMultifactorAuthentication = NewBool(false)
	}

	if s.EnableUserAccessTokens == nil {
		s.EnableUserAccessTokens = NewBool(false)
	}

	if s.GoroutineHealthThreshold == nil {
		s.GoroutineHealthThreshold = NewInt(-1)
	}

	if s.ConnectionSecurity == nil {
		s.ConnectionSecurity = NewString("")
	}

	if s.TLSKeyFile == nil {
		s.TLSKeyFile = NewString(SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE)
	}

	if s.TLSCertFile == nil {
		s.TLSCertFile = NewString(SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE)
	}

	if s.UseLetsEncrypt == nil {
		s.UseLetsEncrypt = NewBool(false)
	}

	if s.LetsEncryptCertificateCacheFile == nil {
		s.LetsEncryptCertificateCacheFile = NewString("./config/letsencrypt.cache")
	}

	if s.ReadTimeout == nil {
		s.ReadTimeout = NewInt(SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT)
	}

	if s.WriteTimeout == nil {
		s.WriteTimeout = NewInt(SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT)
	}

	if s.MaximumLoginAttempts == nil {
		s.MaximumLoginAttempts = NewInt(SERVICE_SETTINGS_DEFAULT_MAX_LOGIN_ATTEMPTS)
	}

	if s.Forward80To443 == nil {
		s.Forward80To443 = NewBool(false)
	}

	if s.SessionIdleTimeoutInMinutes == nil {
		s.SessionIdleTimeoutInMinutes = NewInt(0)
	}

	if s.WebsocketPort == nil {
		s.WebsocketPort = NewInt(80)
	}

	if s.AllowCorsFrom == nil {
		s.AllowCorsFrom = NewString(SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM)
	}

	if s.AllowCookiesForSubdomains == nil {
		s.AllowCookiesForSubdomains = NewBool(false)
	}

	if s.WebserverMode == nil {
		s.WebserverMode = NewString("gzip")
	} else if *s.WebserverMode == "regular" {
		*s.WebserverMode = "gzip"
	}

}

type RateLimitSettings struct {
	Enable           *bool
	PerSec           *int
	MaxBurst         *int
	MemoryStoreSize  *int
	VaryByRemoteAddr *bool
	VaryByUser       *bool
	VaryByHeader     string
}

func (s *RateLimitSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}

	if s.PerSec == nil {
		s.PerSec = NewInt(10)
	}

	if s.MaxBurst == nil {
		s.MaxBurst = NewInt(100)
	}

	if s.MemoryStoreSize == nil {
		s.MemoryStoreSize = NewInt(10000)
	}

	if s.VaryByRemoteAddr == nil {
		s.VaryByRemoteAddr = NewBool(true)
	}

	if s.VaryByUser == nil {
		s.VaryByUser = NewBool(false)
	}
}

func (rls *RateLimitSettings) isValid() *AppError {
	if *rls.MemoryStoreSize <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.rate_mem.app_error", nil, "", http.StatusBadRequest)
	}

	if *rls.PerSec <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.rate_sec.app_error", nil, "", http.StatusBadRequest)
	}

	if *rls.MaxBurst <= 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.max_burst.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type ImportInventorySettings struct {
	Enable            *bool
	IblockId          *int
	MatchFieldName    string
	MatchFieldNum     string
	MatchFieldL       string
	MatchFieldPeriod1 string
	MatchFieldPeriod2 string
}

func (s *ImportInventorySettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}

}

func (rls *ImportInventorySettings) isValid() *AppError {

	return nil
}

type ImportUnitSettings struct {
	Enable            *bool
	IblockId          *int
	MatchFieldName    string
	MatchFieldNum     string
	MatchFieldL       string
	MatchFieldPeriod1 string
	MatchFieldPeriod2 string
}

func (s *ImportUnitSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}

}

func (rls *ImportUnitSettings) isValid() *AppError {

	return nil
}

type JobSettings struct {
	RunJobs      *bool
	RunScheduler *bool
}

func (s *JobSettings) SetDefaults() {
	if s.RunJobs == nil {
		s.RunJobs = NewBool(true)
	}

	if s.RunScheduler == nil {
		s.RunScheduler = NewBool(true)
	}
}

func (o *Config) Sanitize() {

	*o.SqlSettings.DataSource = FAKE_SETTING
	o.SqlSettings.AtRestEncryptKey = FAKE_SETTING

	for i := range o.SqlSettings.DataSourceReplicas {
		o.SqlSettings.DataSourceReplicas[i] = FAKE_SETTING
	}

	for i := range o.SqlSettings.DataSourceSearchReplicas {
		o.SqlSettings.DataSourceSearchReplicas[i] = FAKE_SETTING
	}

}
