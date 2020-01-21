# build
FROM golang:latest AS builder

WORKDIR /build
ADD main.go main.go
ADD cpst cpst
ADD go.mod go.mod
ADD go.sum go.sum
RUN mkdir $GOPATH/src/Copy-Paste && mv * $GOPATH/src/Copy-Paste/
RUN cd $GOPATH/src/Copy-Paste && go mod download && CGO_ENABLED=0 go build -ldflags "-s -w" -o server main.go && mv server /build/

# run
FROM alpine:latest

WORKDIR /root
ADD resources resources
COPY --from=builder /build/server .
EXPOSE 80

CMD ["./server"]
