package main

import "log"

// Config defines expected values to run benchark
type RunList struct {
	reads, writes map[string]string
}

// Result is
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

	// Run configured read benchmarks
	for name, query := range runList.reads {

		time, err := kdb.TimedPost("/api/v1/datapoints/query", query)
		if err != nil {
			log.Printf("Failed read benchmark %s because %s \n", name, err)
			continue
		}

		results["read."+name] = Result{name: name, group: "read", timeMs: time}
	}

	// Run configured write benchmarks
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
