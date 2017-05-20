package main

import (
	// "io/ioutil"

	"net/http/httptest"
	"regexp"
	"testing"
)

func BenchmarkTest0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile0(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile1(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile2(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile3(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile4(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}
