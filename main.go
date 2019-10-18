package main

import (
	"./cpst"
	"flag"
	"fmt"
)

func main() {
	listenAddrParam := flag.String("l", "0.0.0.0:80", "listen address")
	redisHostParam := flag.String("rh", "redis", "redis server host name")
	postgresHostParam := flag.String("ph", "postgres", "postgres server host name")
	flag.Parse()
	listenAddr := *listenAddrParam
	redisHost := *redisHostParam
	postgresHost := *postgresHostParam
	//modify code len and encode chars
	cpst.SetCodeLen(7)
	cpst.SetEncodeChars("苟利国家生死以岂因祸福避趋之")
	server := cpst.NewServer(fmt.Sprintf("%s:6379", redisHost), fmt.Sprintf("postgres://cpst:cpst@%s/cpst?sslmode=disable", postgresHost))
	err := server.Start(listenAddr)
	if err != nil {
		panic(err)
	}
}
