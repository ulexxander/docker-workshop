FROM golang:1.17 AS builder
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
WORKDIR /build
COPY . .
RUN go build main.go

FROM alpine
WORKDIR /docker-workshop
COPY --from=builder /build/main .
CMD ["./main"]
