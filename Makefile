PREFIX=/usr/local
GOBIN=${DESTDIR}${PREFIX}/bin

build:
	go get
	go build

install: build
	GOBIN=${GOBIN} go install -v
	install -Dm 644 config/config.yaml ${DESTDIR}/etc/kairosdb-pref/config.yaml
