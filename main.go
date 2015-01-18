package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

func main() {

	//args
	daemon := flag.Bool("d", false, "run forever")
	flag.Parse()

	// Setup configuration
	viper.SetConfigName("config")

	viper.AddConfigPath("./config/")           //dev config (overrides live)
	viper.AddConfigPath("/etc/kairosdb-pref/") //live config

	viper.ReadInConfig()

	if *daemon == true {

		if !viper.IsSet("frequency") {
			log.Panic("frequency option missing from config. Cannot run daemon")
		}

		fmt.Println("---------------------------------------------------------------------")
		fmt.Printf("Starting Daemon... (benchmark every %v seconds)\n", viper.GetInt("frequency"))
		fmt.Println("---------------------------------------------------------------------")

		var lastRunTs int32
		resultChannel := make(chan []Result)

		//don't let anyone DDoS themselves
		if viper.GetInt("frequency") == 0 {
			log.Panic("benchmark frequency must be greater than 0")
		}

		for /* ever */ {

			if int32(time.Now().Unix()) >= lastRunTs+int32(viper.GetInt("frequency")) {

				log.Println("Staring benchmark... ")

				lastRunTs = int32(time.Now().Unix())

				//always use an up-to-date config
				viper.ReadInConfig()

				//exeute benchmark
				go func(c chan []Result) {
					resultChannel <- run()
				}(resultChannel)
			}

			select {

			//print results when ready
			case result := <-resultChannel:
				printResults(result)

			//check for pending benchmarks every second
			case <-time.After(time.Second):
				continue
			}
		}
	} else {
		printResults(run())
	}
}

func run() []Result {

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
	result := RunBenchmark(&kdb, &runList)

	// Output result
	var datapoints []Datapoint
	for _, result := range result {
		hostname, err := os.Hostname()

		if err != nil {
			hostname = "unknown"
			log.Printf("os.Hostname returned error: %s", err)
		}

		tags := map[string]string{"name": result.name, "group": result.group, "host": hostname}

		if viper.IsSet("logback_tags") {
			for name, value := range viper.GetStringMapString("logback_tags") {
				tags[name] = value
			}
		}

		datapoint := Datapoint{
			Name:      "kairosdb.benchmark.result",
			Timestamp: kdb.MsTime(),
			Value:     result.timeMs,
			Tags:      tags,
		}

		datapoints = append(datapoints, datapoint)
	}

	// flush new datapoints back to kairosdb
	if viper.GetBool("logback") {
		err := kdb.AddDatapoints(datapoints)
		if err != nil {
			log.Printf("logback failed with error: %s", err)
		}
	}

	return result
}

func printResults(result []Result) {
	log.Println("RESULTS:")
	for _, result := range result {
		if result.success {
			log.Printf("[%s] %s success in %d ms", result.group, result.name, result.timeMs)
		} else {
			log.Printf("[%s] %s failed due to: %s", result.group, result.name, result.err)
		}

	}
}
