FROM golang:stretch

WORKDIR /gin-api

COPY gin_server.go .

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go mod init gin_api
RUN CGO_ENABLED=0 GOOS=linux go build

CMD ["./gin_api"]
