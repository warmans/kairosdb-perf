package main

import (
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func main() {

	//setup configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("/var/www/kairosdb-perf/config/")
	viper.ReadInConfig()

	//create kairosdb client
	kdb := Kairosdb{
		client: &http.Client{Timeout: (time.Duration(viper.GetInt("host")) * time.Second)},
		host:   viper.GetString("host"),
	}

	//get benchmarks
	var reads map[string]string
	if viper.IsSet("reads") {
		reads = viper.GetStringMapString("reads")
	}

	var writes map[string]string
	if viper.IsSet("writes") {
		writes = viper.GetStringMapString("writes")
	}

	runList := RunList{
		reads:  reads,
		writes: writes,
	}

	//run benchmarks and return Result
	result := RunBenchmark(kdb, runList)

	//output result
	var datapoints []Datapoint
	for _, result := range result {
		log.Printf("%s (%s) completed in %d ms", result.name, result.group, result.timeMs)

		datapoint := Datapoint{
			Name:      "kairosdb.benchmark.result",
			Timestamp: kdb.MsTime(),
			Value:     result.timeMs,
			Tags:      map[string]string{"name": result.name, "group": result.group},
		}

		datapoints = append(datapoints, datapoint)
	}

	//flush new datapoints back to kairosdb
	if viper.GetBool("logback") {
		log.Print("logging results back to kairosdb")
		kdb.AddDatapoints(datapoints)
	} else {
		log.Print("Discarding result (logback is false)")
	}
}
