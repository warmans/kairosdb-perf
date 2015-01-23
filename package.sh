#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

#set target package format
TARGET=${1-rpm}

#help text
if [ "$TARGET" == "help" ]; then
    echo "Usage ./package.sh [help|rpm|deb|solaris|puppet]"
    exit
fi

#validate fpm is installed
if ! which fpm &>/dev/null; then
    printf "error: Packaging requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"
    exit 1
fi

echo "Building $TARGET package..."

#build in pkg directory
export DESTDIR=pkg

#run install using dest DESTDIR prefix
make install PREFIX=/usr

#clean
if [ -d "dist" ]; then
    rm -f dist/*.$TARGET
else
    mkdir dist
fi

#build RPM
fpm -s dir -p dist -t $TARGET -n kairosdb-perf -v $(cat version) -C $DESTDIR .
