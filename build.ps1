mkdir out -ErrorAction SilentlyContinue
Push-Location src
try {
    go build -v -x -o ../out -ldflags "-w -s" -trimpath
}
finally {
    Pop-Location
}
