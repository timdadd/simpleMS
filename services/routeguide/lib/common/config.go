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

// AppConfig is global object that holds all application level variables.
var AppConfig appConfig

type appConfig struct {
	//	Viper parameters (env or file)
	V *viper.Viper
	// Logging Details
	Log *logrus.Logger

	keyPrefix string // Prefix for key
}

// LoadConfig loads AppConfig from files, command line, environment etc.
func LoadConfig(serviceName string, defaults string, configPaths ...string) (appConfig, error) {
	log := logrus.New()
	AppConfig.Log = log
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

	// Use viper library to load the configuration
	v := viper.New()
	AppConfig.V = v
	v.SetConfigType("yaml")
	v.AutomaticEnv() // Automatically read environment variables

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

	// Note that current directory is already added
	log.Debug("Adding the file paths:", configPaths)
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	configName := serviceName + ".yaml"
	log.Debug("Loading the AppConfig file:", configName)
	v.SetConfigName(configName)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No AppConfig file found")
		} else {
			return AppConfig, fmt.Errorf("failed to read the configuration file: %s:%w", configName, err)
		}
	}

	log.Debug(v.AllSettings())

	return AppConfig, nil
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

// Environment variables take priority
func (c *appConfig) GetBoolKey(key string) bool {
	b, _ := strconv.ParseBool(c.GetStringKey(key))
	return b

}

// The port to listen on
func (c *appConfig) Port() int {
	return c.GetIntKey("port")
}

// The service version
func (c *appConfig) ServiceVersion() string {
	return c.GetStringKey("service_version")
}

// The server address in the format of host:port
func (c *appConfig) ServiceAddress() string {
	return c.GetStringKey("service_addr")
}

// Connection uses TLS if true, else plain TCP
func (c *appConfig) TLS() bool {
	return c.GetBoolKey("tls")
}

// The file containing the CA root cert file
func (c *appConfig) CertFile() string {
	return c.GetStringKey("ca_file")
}

// The server name used to verify the hostname returned by the TLS handshake
func (c *appConfig) HostOverride() string {
	return c.GetStringKey("server_host_override")
}

// The TLS key file
func (c *appConfig) KeyFile() string {
	return c.GetStringKey("key_file")
}

// The address to listen on?
func (c *appConfig) ListenAddress() string {
	return c.GetStringKey("listen_addr")
}
