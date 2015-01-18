# KairosDB Perf

KairosDB is a timeseries database: https://github.com/kairosdb/kairosdb

kairosdb-perf is a tool that can be used to setup an automated performance benchmark so that
real-world performance can be monitored over time.

## Usage

Define read and write queries in the config.yaml (installed /etc/kairosdb-perf/config.yaml).

The programme has two modes of operation; single and daemon. Running with no arguments will execute
the configured benchmarks once and return the results. Running with the `-d` flag will cause the
benchmarks to run every `frequency` seconds (default 60) forever.

An init.d script is installed to manage the daemon which can be run as normal: `service kairosdb-perf start`

Assuming logback is enabled in the config result timings will be logged back to kairosdb
in the kairosdb.benchmark.result metric namespace.

### Notes on Configuration

The config is re-loaded each time the benchmark is run so any changes to config do not require the service
to be restarted.

## Installing from source

Use the Makefile to install from source. Note that you must supply your GOPATH:

```bash
sudo GOPATH=/your/go/path make install
```

## Packaging for distribution

Requires https://github.com/jordansissel/fpm

To create a deployable package (rpm/deb/etc.) use `./package.sh`. The script takes one argument
to specify the target format (e.g. rpm or deb) which will default to rpm if no argument is supplied.

```bash
./package.sh
```
