FROM golang:1.17 AS builder
WORKDIR /build
COPY . .
RUN go build main.go

FROM alpine
WORKDIR /docker-workshop
COPY --from=builder /build/main .
CMD ["./main"]
