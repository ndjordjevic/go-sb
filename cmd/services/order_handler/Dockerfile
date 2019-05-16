#first stage - builder
FROM golang:stretch as builder
ENV GO111MODULE=on
WORKDIR /order_handler
COPY go.mod .
RUN go mod download
COPY order_handler.go .
RUN CGO_ENABLED=0 GOOS=linux go build

#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /order_handler .
EXPOSE 50051
CMD ["./order_handler"]
