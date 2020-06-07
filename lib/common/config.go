package common

import (
	"fmt"
	"strconv"

	//"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

// ServiceConfig is global object that holds all application level variables.
var ServiceConfig appConfig

// Enumerated keys related to the configuration
type key struct {
	ServiceVersion string // The service version
	ServiceAddress string // The server address in the format of host:port
	TLS            string // Connection uses TLS if true, else plain TCP
	CAFile         string // The file containing the CA root cert file
	HostOverride   string // The server name used to verify the hostname returned by the TLS handshake
	ListenAddress  string // The address to listen on
	Port           string // The port to listen on
}

type appConfig struct {
	//	Viper
	V *viper.Viper
	// Logging Details
	Log *logrus.Logger

	keyPrefix string // Prefix for key
	Key       key
}

// LoadConfig loads ServiceConfig from files, command line, environment etc.
func LoadConfig(configName string, defaults string, configPaths ...string) (appConfig, error) {
	log := logrus.New()
	ServiceConfig.Log = log
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{TimestampFormat: time.RFC822}
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

	// Values of keys in the ServiceConfig
	ServiceConfig.Key = key{
		ServiceVersion: "service_version",      // The service version
		ServiceAddress: "service_addr",         // The server address in the format of host:port
		TLS:            "tls",                  // Connection uses TLS if true, else plain TCP
		CAFile:         "ca_file",              // The file containing the CA root cert file
		HostOverride:   "server_host_override", // The server name used to verify the hostname returned by the TLS handshake
		ListenAddress:  "listen_addr",          // The address to listen on
		Port:           "port",                 // The address to listen on
	}

	// Use viper library to load the configuration
	v := viper.New()
	ServiceConfig.V = v
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AutomaticEnv() // Automatically read environment variables

	log.Debug("Loading the defaults")
	v.AddConfigPath(".")
	v.SetConfigName("defaultConfig.yaml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No default ServiceConfig file found")
		} else {
			log.Debug("Error loading defaultConfig.yaml: ", err)
		}
	}

	//	Add any defaults provided by the calling procedure
	v.MergeConfig(strings.NewReader(defaults))

	// Note that current directory is already added
	log.Debug("Adding the file paths:", configPaths)
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	log.Debug("Loading the ServiceConfig file:", configName)
	v.SetConfigName(configName)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No ServiceConfig file found")
		} else {
			return ServiceConfig, fmt.Errorf("failed to read the configuration file: %s:%w", configName, err)
		}
	}

	log.Debug(v.AllSettings())

	return ServiceConfig, nil
}

// Environment variables take priority
func (c *appConfig) KeyPrefix(p string) {
	c.keyPrefix = p
	c.V.SetEnvPrefix(p)
	c.Log.Debug("+ ", c.keyPrefix)
}

// Environment variables take priority
func (c *appConfig) GetStringKey(key string) string {
	//c.V.SetEnvPrefix(c.KeyPrefix)
	s := c.V.GetString(key)
	if s == "" {
		s = c.V.GetString(c.keyPrefix + "." + key)
	}
	c.Log.Debug("+---- ", key, " = ", s)
	return s
}

// Environment variables take priority
func (c *appConfig) GetIntKey(key string) int {
	i, _ := strconv.Atoi(c.GetStringKey(key))
	return i
}
