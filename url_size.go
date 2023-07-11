package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

const TimeOutSeconds = 15

type result struct {
	url  string
	size int64
}

func process(wg *sync.WaitGroup, queue chan result, client http.Client, url string) {
	defer wg.Done()
	res, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request[%s] error[%s]\n", url, err)
		return
	}
	size := res.ContentLength
	if res.Request.Method == "HEAD" || size == -1 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request[%s] error[%s]\n", url, err)
			return
		}
		res.Body.Close()
		size = int64(len(body))
	}

	queue <- result{url: url, size: size}
}

func processUrl(urls []string) []result {
	queue := make(chan result, len(urls))
	var wg sync.WaitGroup
	var responseList []result

	client := http.Client{
		Timeout: TimeOutSeconds * time.Second,
	}
	for _, url := range urls {
		wg.Add(1)
		go process(&wg, queue, client, url)
	}

	wg.Wait()
	close(queue)

	for {
		entry, ok := <-queue
		if !ok {
			break
		}
		responseList = append(responseList, entry)
	}

	sort.SliceStable(responseList, func(j, k int) bool {
		return responseList[j].size < responseList[k].size
	})
	return responseList
}

func main() {
	fmt.Println("start...")

	responseList := processUrl(os.Args[1:])
	for _, v := range responseList {
		fmt.Printf("Url[%s] Response size[%d]\n", v.url, v.size)
	}

	fmt.Println("done!")
}
