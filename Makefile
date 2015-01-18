PREFIX=/usr/local
GOBIN=${DESTDIR}${PREFIX}/bin

build:
	go get
	go build

install: build

	#install binary
	GOBIN=${GOBIN} go install -v

	#install config file
	install -Dm 644 config/config.yaml ${DESTDIR}/etc/kairosdb-pref/config.yaml

	#install init script
	install -Dm 755 init.d/kairosdb-perf ${DESTDIR}/etc/init.d/kairosdb-perf
