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

func BenchmarkTest5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile5(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := []string{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile6(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}
}

func BenchmarkTest7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := map[string]bool{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile7(&FileInfo{
			filePath:       filePath,
			seenBrowsers:   seenBrowsers,
			r:              r,
			w:              w,
			uniqueBrowsers: uniqueBrowsers,
		})
	}
}

func BenchmarkTest8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filePath := "data/logs0.txt"
		seenBrowsers := map[string]bool{}
		uniqueBrowsers := 0
		w := httptest.NewRecorder()
		r := regexp.MustCompile("@")
		SearchFile8(&FileInfo{
			filePath:       filePath,
			seenBrowsers:   seenBrowsers,
			r:              r,
			w:              w,
			uniqueBrowsers: uniqueBrowsers,
		})
	}
}
