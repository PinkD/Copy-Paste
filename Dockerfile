# build
FROM golang:latest AS builder

WORKDIR /build
COPY main.go .
ADD cpst cpst

RUN cd cpst && go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o server main.go

# run
FROM alpine:latest

WORKDIR /root
ADD resources resources
COPY --from=builder /build/server .
EXPOSE 80

CMD ["./server"]
