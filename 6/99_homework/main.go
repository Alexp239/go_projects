package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	orderAsc = iota
	orderDesc
)

var (
	errTest = errors.New("testing")
	client  = &http.Client{Timeout: time.Second}
)

type User struct {
	Id     int
	Name   string
	Age    int
	About  string
	Gender string
}

type SearchResponse struct {
	Users    []User
	NextPage bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func doSearch(searcherURL string, limit int, offset int, query string, orderField string, orderBy int) (*SearchResponse, error) {

	searcherParams := url.Values{}

	if limit < 0 {
		return nil, fmt.Errorf("limit must be > 0")
	}
	if limit > 25 {
		limit = 25
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset must be > 0")
	}

	//нужно для получения следующей записи, на основе которой мы скажем - можно показать переключатель следующей страницы или нет
	limit++

	searcherParams.Add("limit", strconv.Itoa(limit))
	searcherParams.Add("offset", strconv.Itoa(offset))
	searcherParams.Add("query", query)
	searcherParams.Add("order_field", orderField)
	searcherParams.Add("order_by", strconv.Itoa(orderBy))

	limit--

	req, err := http.NewRequest("GET", searcherURL+"?"+searcherParams.Encode(), nil)
	// Была бы ошибка, если параметры у http.NewRequest указаны неверно
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, fmt.Errorf("timeout for %s", searcherParams.Encode())
		}
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	result := SearchResponse{}

	data := []User{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	result.Users = data[0:min(len(data), limit)]
	if len(data) > limit {
		result.NextPage = true
	}

	// fmt.Printf("%+v\n", data)

	return &result, err
}
