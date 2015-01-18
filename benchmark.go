package main

// RunList defines the queries to be run by the benchmark [name]query
type RunList struct {
	reads, writes map[string]string
}

// Result instances hold results of benchmarks
type Result struct {
	name, group string
	success     bool
	timeMs      int64
	err         error
}

// RunBenchmark runs all benchmarks
func RunBenchmark(kdb *Kairosdb, runList *RunList) []Result {

	var results []Result

	results = doVersion(kdb, results)
	results = doReads(runList.reads, kdb, results)
	results = doWrites(runList.writes, kdb, results)

	return results
}

// doVersion: Benchmark version endpoint to get a baseline for communicating with kairos
func doVersion(kdb *Kairosdb, results []Result) []Result {
	timeMs, err := kdb.TimedGet("/api/v1/version")

	result := Result{name: "version", group: "version", timeMs: timeMs}

	if err != nil {
		result.err = err
		result.success = false
	} else {
		result.success = true
	}

	return append(results, result)
}

//doReads: Perform read benchmark
func doReads(reads map[string]string, kdb *Kairosdb, results []Result) []Result {

	for name, query := range reads {

		timeMs, err := kdb.TimedPost("/api/v1/datapoints/query", query)

		result := Result{name: name, group: "read", timeMs: timeMs}

		if err != nil {
			result.err = err
			result.success = false
		} else {
			result.success = true
		}

		results = append(results, result)
	}

	return results
}

//doWrites: Perform write benchmark
func doWrites(writes map[string]string, kdb *Kairosdb, results []Result) []Result {

	for name, query := range writes {

		timeMs, err := kdb.TimedPost("/api/v1/datapoints", query)

		result := Result{name: name, group: "write", timeMs: timeMs}

		if err != nil {
			result.err = err
			result.success = false
		} else {
			result.success = true
		}

		results = append(results, result)
	}

	return results
}
