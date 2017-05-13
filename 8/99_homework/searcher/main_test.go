package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := http.Client{}
		resp, err := c.Get("http://127.0.0.1:8081/")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
	}
}
