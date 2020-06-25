package common

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"strconv"
	"sync"

	//"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

//// App is global variable that holds all common application level variables.
var (
	App                 AppConfig
	ErrNoConfigSettings = errors.New("Didn't find any settings")
)

//// AppConfig is global structure that holds all common application level variables.
type AppConfig struct {
	//	Viper parameters (env or file)
	V *viper.Viper
	// Logging Details
	Log          *logrus.Logger
	Mutex        *sync.Mutex
	Ctx          context.Context
	Platform     PlatformDetails
	CanaryColour string
	// grpc Service connections
	SvcConn map[string]*grpc.ClientConn
	// Prefix for key used to find config/environment variables
	keyPrefix string
}

type PlatformDetails struct {
	Url      string
	Provider string
}

// LoadConfig loads AppConfig from files, command line, environment etc.
func LoadConfig(serviceName string, defaults string, configPaths ...string) (AppConfig, error) {
	App.Ctx = context.Background()
	log := logrus.New()
	App.Log = log
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{TimestampFormat: time.RFC822}
	log.WithField("prefix", serviceName)
	//log.Formatter = &logrus.JSONFormatter{
	//	FieldMap: logrus.FieldMap{
	//		logrus.FieldKeyTime:  "timestamp",
	//		logrus.FieldKeyLevel: "severity",
	//		logrus.FieldKeyMsg:   "message",
	//	},
	//	TimestampFormat: time.RFC3339Nano,
	//}
	log.Out = os.Stdout

	log.Debug("Loading the configuration")

	// Use viper library to load the configuration
	v := viper.New()
	App.V = v
	v.SetConfigType("yaml")
	v.AutomaticEnv() // Automatically read environment variables

	log.Debug("Adding the file paths:", configPaths)
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	log.Debug("Loading the defaults")
	v.AddConfigPath("./cfg") // Default directory for config files
	v.SetConfigName("defaultConfig.yaml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No default AppConfig file found")
		} else {
			log.Debug("Error loading defaultConfig.yaml: ", err)
		}
	}

	//	Add any defaults provided by the calling procedure
	v.MergeConfig(strings.NewReader(defaults))

	configName := serviceName + ".yaml"
	log.Debug("Loading the AppConfig file:", configName)
	v.SetConfigName(configName)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No AppConfig file found")
		} else {
			return App, fmt.Errorf("failed to read the configuration file: %s:%w", configName, err)
		}
	}

	settings := v.AllSettings()
	log.Debug(settings)
	if len(settings) == 0 {
		return App, ErrNoConfigSettings
	}

	App.Mutex = &sync.Mutex{}

	App.SvcConn = make(map[string]*grpc.ClientConn)

	//get env and render correct platform banner.
	var env = App.GetStringKey("ENV_PLATFORM")
	platform := PlatformDetails{}
	platform.setPlatformDetails(strings.ToLower(env))
	App.Platform = platform

	App.CanaryColour = v.GetString("CANARY_COLOUR") // Shown on web pages

	return App, nil
}

// Environment variables take priority
func (c *AppConfig) KeyPrefix(p string) {
	c.keyPrefix = p
	c.V.SetEnvPrefix(p)
	c.Log.Debug("+ ", c.keyPrefix)
}

// Environment variables take priority
func (c *AppConfig) GetStringKey(key string) string {
	//c.V.SetEnvPrefix(c.KeyPrefix)
	s := c.V.GetString(key)
	if s == "" {
		s = c.V.GetString(c.keyPrefix + "." + key)
	}
	c.Log.Debug("+---- ", key, " = ", s)
	return s
}

// Environment variables take priority
func (c *AppConfig) GetIntKey(key string) int {
	i, _ := strconv.Atoi(c.GetStringKey(key))
	return i
}

// Environment variables take priority
func (c *AppConfig) GetBoolKey(key string) bool {
	b, _ := strconv.ParseBool(c.GetStringKey(key))
	return b

}

// The port to listen on
func (c *AppConfig) Port() int {
	return c.GetIntKey("port")
}

// The service version
func (c *AppConfig) ServiceVersion() string {
	return c.GetStringKey("service_version")
}

// The server address in the format of host:port
func (c *AppConfig) ServiceAddress() string {
	return c.GetStringKey("service_addr")
}

// Connection uses TLS if true, else plain TCP
func (c *AppConfig) TLS() bool {
	return c.GetBoolKey("tls")
}

// The file containing the CA root cert file
func (c *AppConfig) CertFile() string {
	return c.GetStringKey("ca_file")
}

// The server name used to verify the hostname returned by the TLS handshake
func (c *AppConfig) HostOverride() string {
	return c.GetStringKey("server_host_override")
}

// The TLS key file
func (c *AppConfig) KeyFile() string {
	return c.GetStringKey("key_file")
}

// The address to listen on?
func (c *AppConfig) ListenAddress() string {
	return c.GetStringKey("listen_addr")
}

func (p *PlatformDetails) setPlatformDetails(env string) {
	switch e := strings.ToLower(env); e {
	case "":
		p.Provider = "Test"
		p.Url = ""
	case "alibaba":
		p.Provider = "Alibaba"
		p.Url = "https://us.alibabacloud.com/"
	case "local":
		p.Provider = "local pc"
		p.Url = "http://localhost"
	case "aws":
		p.Provider = "AWS"
		p.Url = "https://aws.amazon.com/"
	case "onprem":
		p.Provider = "On-Premises"
		p.Url = "http://www"
	case "azure":
		p.Provider = "Azure"
		p.Url = "https://azure.microsoft.com/en-gb/"
	default:
		p.Provider = "Google Cloud"
		p.Url = "https://cloud.google.com/"
	}
}
