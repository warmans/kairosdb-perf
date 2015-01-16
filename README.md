Install
-------

```bash
make
sudo make install
```

Package
-------

```bash
export DESTDIR=pkg
make
make install PREFIX=/usr
fpm -s dir -t rpm -n kairosdb-pref -v 0.0.1 $DESTDIR
```
