package main

import "time"

type config struct {
	ApplicationName    string `env:"NAME" default:"go-fibo"`
	ApplicationVersion string `env:"VERSION" default:"v0.0.14"`
	InstanceGUID       string `env:"CF_INSTANCE_GUID"`

	LogLevel   string `env:"LOGLEVEL" default:"INFO"`
	LogAs      string `env:"LOGAS" default:"text"`
	DebugLevel int    `env:"DEBUGLEVEL" default:"0"`
	LogDetails bool   `env:"LOGDETAILS" default:"true"`

	FiboMin int `env:"MIN" default:"16"`
	FiboMax int `env:"MAX" default:"48"`

	Timezone string `env:"TIMEZONE" default:"Local"`

	InstancePort    string `env:"PORT" required:"true"`
	Hostname        string `env:"HOSTNAME"`
	InstanceAddress string `env:"CF_INSTANCE_ADDR"`

	ShutdownWaitTime  time.Duration `env:"SHUTDOWNWAITTIME" default:"480s"`
	MetricsUpdateTime time.Duration `env:"METRICSUPDATETIME" default:"30s"`
}
