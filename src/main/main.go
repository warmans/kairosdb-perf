package main

import (
	"time"

	"github.com/spf13/viper"
)

func main() {

	//setup configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("/var/www/kairosdb-perf/config/")
	viper.ReadInConfig()

	var reads map[string]string
	if viper.IsSet("reads") {
		reads = viper.GetStringMapString("reads")
	}

	var writes map[string]string
	if viper.IsSet("writes") {
		writes = viper.GetStringMapString("writes")
	}

	RunBenchmark(Config{
		host:    viper.GetString("host"),
		timeout: (time.Duration(viper.GetInt("host")) * time.Second),
		reads:   reads,
		writes:  writes,
	})
}
