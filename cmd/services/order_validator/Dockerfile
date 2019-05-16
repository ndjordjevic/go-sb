#first stage - builder
FROM golang:stretch as builder
ENV GO111MODULE=on
WORKDIR /order_validator
COPY go.mod .
RUN go mod download
COPY order_validator.go .
RUN CGO_ENABLED=0 GOOS=linux go build

#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /order_validator .
CMD ["./order_validator"]
