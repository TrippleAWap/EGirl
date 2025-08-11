pushd src
go build -o ../out -ldflags="-w -s" -trimpath
popd
