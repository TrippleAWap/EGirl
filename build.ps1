mkdir out -ErrorAction SilentlyContinue
Push-Location src
try {
    go build -o ../out -ldflags "-w -s" -trimpath
}
finally {
    Pop-Location
}
