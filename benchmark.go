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

// Benchmark class
type Benchmark struct {
	kdb     *Kairosdb
	results []Result
}

// Execute runs all benchmarks in supplied RunList
func (benchmark *Benchmark) Execute(runList *RunList) []Result {

	benchmark.doVersion()
	benchmark.doPostBenchmark(runList.reads, "read", "/api/v1/datapoints/query")
	benchmark.doPostBenchmark(runList.writes, "write", "/api/v1/datapoints")

	return benchmark.results
}

// doVersion: Benchmark version endpoint to get a baseline for communicating with kairos
func (benchmark *Benchmark) doVersion() {
	timeMs, err := benchmark.kdb.TimedGet("/api/v1/version")

	result := Result{name: "version", group: "version", timeMs: timeMs}

	if err != nil {
		result.err = err
		result.success = false
	} else {
		result.success = true
	}

	benchmark.results = append(benchmark.results, result)
}

// doPostBenchmark: runs a POST request against a resource and appends result to results
func (benchmark *Benchmark) doPostBenchmark(queries map[string]string, group string, path string) {

	for name, query := range queries {

		timeMs, err := benchmark.kdb.TimedPost(path, query)

		result := Result{name: name, group: group, timeMs: timeMs}

		if err != nil {
			result.err = err
			result.success = false
		} else {
			result.success = true
		}

		benchmark.results = append(benchmark.results, result)
	}

}
