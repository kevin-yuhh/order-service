package config

import (
	"github.com/TRON-US/soter-order-service/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	ConnectionUriMaster string
	ConnectionUriSlave  string
	MaxIdleConn         int
	MaxOpenConn         int
	MaxLifetime         int
}

type ServerConfiguration struct {
	Host string
	Port int
}

type LoggerConfiguration struct {
	Level  int8
	Output string
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	ScriptId int
	Logger   LoggerConfiguration
	EvChan   chan bool
}

// Get config struct from file.
func NewConfiguration(configName string, configPath string) (*Configuration, error) {
	// New struct of configure.
	c := &Configuration{EvChan: make(chan bool)}

	// Set config name.
	viper.SetConfigName(configName)
	// Set config path.
	viper.AddConfigPath(configPath)

	// Read config.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config struct.
	err := viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	// Watcher file changed.
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Logger.Info("The configuration file has been modified: ", e.Name)
		err := viper.Unmarshal(c)
		if err != nil {
			logger.Logger.Errorf("Unable to decode into struct, reasons: [%v]", err)
			return
		}
		c.EvChan <- true
	})
	return c, nil
}
