# brew install FiloSottile/musl-cross/musl-cross
# apt-get install -y musl
CC=x86_64-linux-musl-gcc CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app-linux -ldflags="-s -w" main.go
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o app-mac -ldflags="-s -w" main.go

