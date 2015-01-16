export DESTDIR=pkg
make
make install PREFIX=/usr

rm -f *.rpm
fpm -s dir -t rpm -n kairosdb-pref -v 0.0.1 -C $DESTDIR .
