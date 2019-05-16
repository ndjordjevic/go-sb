#first stage - builder
FROM golang:stretch as builder
ENV GO111MODULE=on
WORKDIR /price_handler
COPY go.mod .
RUN go mod download
COPY price_handler.go .
RUN CGO_ENABLED=0 GOOS=linux go build

#second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /price_handler .
EXPOSE 50071
CMD ["./price_handler"]
