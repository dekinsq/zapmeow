FROM golang:1.23-alpine

RUN apk add --no-cache gcc musl-dev
RUN apk add mailcap

WORKDIR /app

COPY . .

RUN ls -la

RUN go mod download

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN go build -o server main.go

EXPOSE 3000

CMD ["./server"]
