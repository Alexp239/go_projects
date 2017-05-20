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

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	r := regexp.MustCompile("@")

	files, _ := ioutil.ReadDir(logsPath)

	for _, f := range files {
		fmt.Fprintln(w, f.Name())
		filePath := logsPath + f.Name()
		SearchFile(filePath, w, &seenBrowsers, &uniqueBrowsers, r)
	}

}

func SearchFile(filePath string, w http.ResponseWriter, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	// fileContents, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// lines := strings.Split(string(fileContents), "\n")

	fmt.Fprintln(w, "found users:\n")
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		user := structs.User{}
		// err := json.Unmarshal([]byte(line), &user)
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: []byte(line)})
		if err != nil {
			// fmt.Println("cant unmarshal json: ", f.Name(), line, err)
			continue
		}
		parseUser(user, seenBrowsers, uniqueBrowsers, r, i, w)
		i++
	}

	fmt.Fprintln(w, "Total unique browsers", *uniqueBrowsers)
}

func parseUser(user structs.User, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp, i int, w http.ResponseWriter) {
	isAndroid := false
	isMSIE := false

	browsers, ok := user.Browsers.([]interface{})
	if !ok {
		// log.Println("cant cast browsers")
		return
	}

	for _, browserRaw := range browsers {
		browser, ok := browserRaw.(string)
		if !ok {
			// log.Println("cant cast browser to string")
			continue
		}

		if matchAndroid(browser) != -1 {
			isAndroid = true
			notSeenBefore := true
			for _, item := range *seenBrowsers {
				if item == browser {
					notSeenBefore = false
				}
			}
			if notSeenBefore {
				// log.Printf("New browser: %s, first seen: %s", browser, user["name"])
				*seenBrowsers = append(*seenBrowsers, browser)
				*uniqueBrowsers++
			}
		}
	}

	for _, browserRaw := range browsers {
		browser, ok := browserRaw.(string)
		if !ok {
			// log.Println("cant cast browser to string")
			continue
		}
		if matchMSIE(browser) != -1 {
			isMSIE = true
			notSeenBefore := true
			for _, item := range *seenBrowsers {
				if item == browser {
					notSeenBefore = false
				}
			}
			if notSeenBefore {
				// log.Printf("New browser: %s, first seen: %s", browser, user["name"])
				*seenBrowsers = append(*seenBrowsers, browser)
				*uniqueBrowsers++
			}
		}
	}

	if !(isAndroid && isMSIE) {
		return
	}

	// log.Println("Android and MSIE user:", user["name"], user["email"])
	email := r.ReplaceAllString(user.Email.(string), " [at] ")
	fmt.Fprintln(w, "[%d] %s <%s>\n", i, user.Name, email)
}
