# shadowsocks-multiuser
## Compile
```
go get -u -v github.com/NetchX/shadowsocks-multiuser
cd $GOPATH/src/github.com/NetchX/shadowsocks-multiuser
go build
```

## Use
```
Usage of shadowsocks-multiuser:
  -dbhost string
        database host (default "localhost")
  -dbname string
        database name (default "sspanel")
  -dbpass string
        database pass (default "123456")
  -dbport int
        database port (default 3306)
  -dbuser string
        database user (default "root")
  -listcipher
        list cipher
  -nodeid int
        node id (default -1)
```