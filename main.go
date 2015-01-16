package main

import (
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// RunList defines the queries to be run by the benchmark [name]query
type RunList struct {
	reads, writes map[string]string
}

// Result instances hold results of benchmarks
type Result struct {
	name, group string
	timeMs      int64
}

// RunBenchmark runs all benchmarks
func RunBenchmark(kdb Kairosdb, runList RunList) map[string]Result {

	log.Printf("Starting benchmark on kairosdb host %s...\n", kdb.host)

	results := make(map[string]Result)

	// Benchmark version endpoint to get a baseline for communicating with kairos
	getTime, err := kdb.TimedGet("/api/v1/version")
	if err != nil {
		log.Printf("Failed read benchmark version because %s \n", err)
	}
	results["version"] = Result{name: "version", group: "read", timeMs: getTime}

	// Run read benchmarks
	for name, query := range runList.reads {

		time, err := kdb.TimedPost("/api/v1/datapoints/query", query)
		if err != nil {
			log.Printf("Failed read benchmark %s because %s \n", name, err)
			continue
		}

		results["read."+name] = Result{name: name, group: "read", timeMs: time}
	}

	// Run write benchmarks
	for name, datapoints := range runList.writes {

		time, err := kdb.TimedPost("/api/v1/datapoints", datapoints)
		if err != nil {
			log.Printf("Failed write benchmark %s because %s \n", name, err)
			continue
		}

		results["write."+name] = Result{name: name, group: "write", timeMs: time}
	}

	return results
}

func main() {

	// Setup configuration
	viper.SetConfigName("config")

	// todo: fix config path
	viper.AddConfigPath("./config/")
	viper.ReadInConfig()

	// Kairosdb client
	kdb := Kairosdb{
		client: &http.Client{Timeout: (time.Duration(viper.GetInt("host")) * time.Second)},
		host:   viper.GetString("host"),
	}

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

	// Run benchmarks and return Result
	result := RunBenchmark(kdb, runList)

	// Output result
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

	// flush new datapoints back to kairosdb
	if viper.GetBool("logback") {
		log.Print("logging results back to kairosdb")
		err := kdb.AddDatapoints(datapoints)
		if err != nil {
			log.Print(err)
		}
	} else {
		log.Print("Discarding result (logback is false)")
	}
}
