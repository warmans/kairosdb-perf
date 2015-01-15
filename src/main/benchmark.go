package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
)

// Config defines expected values to run benchark
type Config struct {
	host          string
	timeout       time.Duration
	reads, writes map[string]string
}

// Result is
type Result struct {
	name, group string
	timeMs      int64
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func timedGet(uri string) int64 {

	start := makeTimestamp()

	res, err := http.Get(uri)
	defer res.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	if err != nil {
		log.Panic(err)
	}

	end := makeTimestamp()

	return (end - start)
}

func timedPost(uri string, body string, timeout time.Duration) (int64, error) {

	start := makeTimestamp()

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	//create HTTP client
	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, errors.New("Response was non-200: " + resp.Status)
	}

	end := makeTimestamp()

	return (end - start), nil
}

// RunBenchmark runs all benchmarks
func RunBenchmark(config Config) map[string]Result {

	log.Printf("Starting benchmark on kairosdb host %s...\n", config.host)

	readResults := make(map[string]Result)

	// Benchmark version endpoint
	readResults["version"] = Result{name: "version", group: "read", timeMs: timedGet(config.host + "/api/v1/version")}

	// Run configured read benchmarks
	for name, query := range config.reads {

		time, err := timedPost(config.host+"/api/v1/datapoints/query", query, config.timeout)
		if err != nil {
			log.Printf("Failed read benchmark %s because %s \n", name, err)
			continue
		}

		readResults[name] = Result{name: name, group: "read", timeMs: time}
	}

	return readResults
}
