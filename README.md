# gifseekr
Find and share animated GIFs.

## local dev (osx)

If you don't have Go installed, follow the instructions
[here](https://golang.org/doc/install).

Get the code:
```console
mkdir -p $GOPATH/github.com && cd $GOPATH/github.com/
git clone github.com/dustinrohde/gifseekr
cd dustinrohde/gifseekr/
```

[Install Go QT bindings](https://github.com/therecipe/qt/wiki/Installation):
```console
# install xcode
xcode-select --install

# follow the instructions to install bindings on osx in go module mode

# download packages
go get -u -v github.com/therecipe/qt/cmd/qtsetup
go get -u -v github.com/therecipe/qt/cmd/...
go mod vendor
git clone https://github.com/therecipe/env_darwin_amd64_513.git vendor/github.com/therecipe/env_darwin_amd64_513

# these two additional steps were necessary so that qtsetup wouldn't fail
rm -rf vendor/github.com/therecipe/qt
git clone https://github.com/therecipe/qt.git vendor/github.com/therecipe/qt

# run qtsetup
$(go env GOPATH)/bin/qtsetup -test=false
```


Hack away!
```console
# Run tests
go test

# Build the app
go build

# Run the app
go run cmd/gifseekr/main.go
```
