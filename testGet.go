package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pinkd.moe/x/copy-paste-text/cpst"
	"strconv"
	"sync"
	"time"
)

func main() {
	count, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		panic(err)
	}
	enableLog := false
	threadCount := 2048
	requestChan := make(chan uint64, threadCount)
	reuseChan := make(chan uint64, threadCount/2)
	wg := sync.WaitGroup{}
	okChan := make(chan interface{})
	start := time.Now()
	go func() {
		for i := 0; i < int(count); i++ {
			for ; len(reuseChan) != 0; { //check reuse
				log.Printf("reuse: %d\n", len(reuseChan))
				time.Sleep(time.Millisecond * 100)
			}
			requestChan <- uint64(i)
		}
		fmt.Println("content is over")
		okChan <- nil
	}()
	client := http.Client{
		//Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        threadCount,
			MaxIdleConnsPerHost: threadCount,
		},
	}
	for i := 0; i < cap(requestChan); i++ {
		go func() {
			for {
				var line uint64
				if len(reuseChan) > 0 { //reuse
					//log.Println("reuse")
					line = <-reuseChan
				} else {
					line = <-requestChan
				}
				wg.Add(1)
				start := time.Now()
				lineChar := cpst.NumberToChar(line)
				response, err := client.Get(fmt.Sprintf("http://192.168.1.1:8080/%s", lineChar))
				if err != nil { //request fail, reuse
					reuseChan <- line
					log.Println(err)
					wg.Done()
					continue
				}
				if response.Body != nil {
					_, _ = io.Copy(ioutil.Discard, response.Body)
					_ = response.Body.Close()
					if err != nil {
						panic(err)
					}
					elapsed := time.Since(start)
					if enableLog {
						fmt.Printf("%v\telapsed, getting %s\n", elapsed, lineChar)
					}
				}
				wg.Done()
			}
		}()
	}
	<-okChan
	wg.Wait()
	end := time.Since(start)
	fmt.Printf("%v elapsed, %d lines got, %.2f op/s", end, count, float64(count)/end.Seconds())
}
