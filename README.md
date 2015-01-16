# KairosDB Perf

KairosDB is a timeseries database: https://github.com/kairosdb/kairosdb

kairosdb-perf is a tool that can be used to setup an automated performance benchmark so that
real-world performance can be monitored over time.

Usage
---------

Define read and write queries in the config.yaml (installed /etc/kairosdb-perf/config.yaml)
then run the binary `kairosdb-perf`. It is expected this will be run on a cron or simialar.

Assuming logback is enabled in the config result timings will be logged back to kairosdb 
in the kairosdb.benchmark.result metric namespace.


Installing from source
----------

Use the Makefile to install from source:

```bash
make
sudo make install
```

Packaging for distribuition
----------
Requires https://github.com/jordansissel/fpm

To create a deployable package (rpm/deb/etc.) use `./package.sh`. The script takes one argument 
to specify the target format (e.g. rpm or deb) which will default to rpm if no argument is supplied.

```bash
./package.sh
```
