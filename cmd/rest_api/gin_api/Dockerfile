#first stage - builder
FROM golang:stretch as builder
ENV GO111MODULE=on
WORKDIR /gin_api
COPY go.mod .
RUN go mod download
COPY gin_server.go .
RUN CGO_ENABLED=0 GOOS=linux go build

#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /gin_api .
EXPOSE 8010
CMD ["./gin_api"]
