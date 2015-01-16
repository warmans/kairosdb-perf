DESTDIR=/usr/local
GOBIN=${DESTDIR}/bin

build:
	go get
	go build

install: build
	GOBIN=${GOBIN} go install -v
	install -Dm 644 config/config.yaml ${DESTDIR}/etc/kairosdb-pref/config.yaml

