package main

import (
	"log"
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

	result := RunBenchmark(Config{
		host:    viper.GetString("host"),
		timeout: (time.Duration(viper.GetInt("host")) * time.Second),
		reads:   reads,
		writes:  writes,
	})

	for _, result := range result {
		log.Printf("%s (%s) completed in %d", result.name, result.group, result.timeMs)
	}
}
