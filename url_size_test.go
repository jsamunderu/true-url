package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_UrlSize(t *testing.T) {
	const RESP_SIZE int = 102400
	var urls []string
	for i := 0; i < 20; i++ {
		urls = append(urls, fmt.Sprintf("https://httpbin.org/bytes/%d", rand.Intn(RESP_SIZE)))
	}

	responseList := processUrl(urls)
	for i := 1; i < len(responseList); i++ {
		if responseList[i].size < responseList[i-1].size {
			t.Error(fmt.Sprintf("Index[%d]Result{Url: %s, size%d} IS GREATER THAN Index[%d]Result{Url: %s, size%d} IS GREATER THAN Index", (i - 1), responseList[i-1].url, responseList[i-1].size, i, responseList[i].url, responseList[i].size))
		}
	}
}
