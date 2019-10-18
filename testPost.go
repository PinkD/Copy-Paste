package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	enableLog := false
	threadCount := 2048
	requestChan := make(chan string, threadCount)
	reuseChan := make(chan string, threadCount/2)
	wg := sync.WaitGroup{}
	okChan := make(chan interface{})
	start := time.Now()
	lineCount := 0
	go func() {
		content := strings.Replace(string(data), "\r", "\n", 0)
		content = strings.Replace(content, "\n\n", "\n", 0)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			for ; len(reuseChan) != 0; { //check reuse
				log.Printf("reuse: %d\n", len(reuseChan))
				time.Sleep(time.Millisecond * 100)
			}
			if len(strings.TrimSpace(line)) != 0 {
				lineCount++
				requestChan <- line
			}
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
				var line string
				if len(reuseChan) > 0 { //reuse
					//log.Println("reuse")
					line = <-reuseChan
				} else {
					line = <-requestChan
				}
				wg.Add(1)
				start := time.Now()
				//response, err := client.PostForm("http://server/new", url.Values{"content": {line}})
				response, err := client.PostForm("http://192.168.1.1:8080/new", url.Values{"content": {line}})
				if err != nil { //request fail, reuse
					reuseChan <- line
					log.Println(err)
					wg.Done()
					continue
				}
				if response.Body != nil {
					result, err := ioutil.ReadAll(response.Body)
					_ = response.Body.Close()
					if err != nil {
						panic(err)
					}
					elapsed := time.Since(start)
					if enableLog {
						fmt.Printf("%v\telapsed, result is %s", elapsed, string(result))
					}
				}
				wg.Done()
			}
		}()
	}
	<-okChan
	wg.Wait()
	end := time.Since(start)
	fmt.Printf("%v elapsed, %d lines sent, %.2f op/s", end, lineCount, float64(lineCount)/end.Seconds())

}
