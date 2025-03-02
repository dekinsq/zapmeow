CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-linux -ldflags="-s -w" main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o app-mac -ldflags="-s -w" main.go

