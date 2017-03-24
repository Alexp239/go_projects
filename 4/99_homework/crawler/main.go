package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func getBaseTag(body string) string {
	ar := strings.Split(body, "<base href=\"")
	if len(ar) == 1 {
		return ""
	}
	res := strings.Split(ar[1], "\"")[0]
	if res[len(res)-1] != '/' {
		res = ""
	}
	return res
}

func crawler(host string, path string, visited *map[string]bool) (res []string) {
	c := http.Client{}
	resp, err := c.Get(host + path)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return res
	}
	defer resp.Body.Close()
	real_path := resp.Request.URL.Path
	(*visited)[real_path] = true
	res = append(res, real_path)
	body, _ := ioutil.ReadAll(resp.Body)
	body_str := string(body)
	base := getBaseTag(body_str)
	ar := strings.Split(body_str, "<a href=\"")
	for _, s := range ar[1:] {
		next_path := base + strings.Split(s, "\">")[0]
		if _, ok := (*visited)[next_path]; !ok {
			res = append(res, crawler(host, next_path, visited)...)
		}
	}
	return res
}

func Crawl(host string) []string {
	visited := make(map[string]bool)
	return crawler(host, "/", &visited)
}
