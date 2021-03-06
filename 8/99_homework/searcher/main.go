package main

import (
	// "encoding/json"
	"bufio"
	"fmt"
	jlexer "github.com/mailru/easyjson/jlexer"
	"io/ioutil"
	"log"
	"msu-go-11/8/99_homework/searcher/structs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	// "strings"
)

const logsPath = "./data/"

type FileInfo struct {
	filePath       string
	w              http.ResponseWriter
	seenBrowsers   map[string]bool
	uniqueBrowsers int
	r              *regexp.Regexp
}

func main() {
	http.Handle("/", http.HandlerFunc(handleFunc))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleFunc(w http.ResponseWriter, req *http.Request) {

	if req.URL.String() == "/favicon.ico" {
		return
	}

	seenBrowsers := map[string]bool{}
	uniqueBrowsers := 0

	r := regexp.MustCompile("@")

	files, _ := ioutil.ReadDir(logsPath)

	for _, f := range files {
		fmt.Fprintln(w, f.Name())
		filePath := logsPath + f.Name()
		SearchFile(&FileInfo{
			seenBrowsers:   seenBrowsers,
			filePath:       filePath,
			w:              w,
			r:              r,
			uniqueBrowsers: uniqueBrowsers,
		})
	}

}

func SearchFile(fileInf *FileInfo) {
	file, err := os.Open(fileInf.filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	fmt.Fprintln(fileInf.w, "found users:\n")
	i := 0
	for scanner.Scan() {
		user := structs.User1{}
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: scanner.Bytes()})
		if err != nil {
			// fmt.Println("cant unmarshal json: ", f.Name(), line, err)
			continue
		}
		parseUser(user, fileInf, i)
		i++
	}

	fmt.Fprintln(fileInf.w, "Total unique browsers", fileInf.uniqueBrowsers)
}

func parseUser(user structs.User1, fileInf *FileInfo, i int) {
	isAndroid := false
	isMSIE := false

	browsers := user.Browsers

	for _, browserRaw := range browsers {
		browser := browserRaw

		matchA := strings.Contains(browser, "Android")
		matchM := strings.Contains(browser, "MSIE")

		isAndroid = matchA
		isMSIE = matchM

		if isAndroid || isMSIE {
			_, ok := fileInf.seenBrowsers[browser]

			if !ok {
				// log.Printf("New browser: %s, first seen: %s", browser, user.Name)
				fileInf.seenBrowsers[browser] = true
				fileInf.uniqueBrowsers++
			}
		}
	}

	if !(isAndroid && isMSIE) {
		return
	}

	// log.Println("Android and MSIE user:", user.Name, user.Email)
	email := fileInf.r.ReplaceAllString(user.Email, " [at] ")
	fmt.Fprintln(fileInf.w, "[%d] %s <%s>\n", i, user.Name, email)
}
