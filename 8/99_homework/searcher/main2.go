package main

import (
	// "encoding/json"
	"fmt"
	jlexer "github.com/mailru/easyjson/jlexer"
	"io/ioutil"
	"log"
	"msu-go-11/8/99_homework/searcher/structs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strings"
)

func SearchFile2(filePath string, w http.ResponseWriter, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp) {
	foundUsers := ""

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(fileContents), "\n")

	users := make([]structs.User, 0)
	for _, line := range lines {
		user := structs.User{}
		// err := json.Unmarshal([]byte(line), &user)
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: []byte(line)})
		if err != nil {
			// fmt.Println("cant unmarshal json: ", f.Name(), line, err)
			continue
		}
		users = append(users, user)
	}

	for i, user := range users {
		parseUser2(user, seenBrowsers, uniqueBrowsers, r, &foundUsers, i)
	}

	fmt.Fprintln(w, "found users:\n"+foundUsers)
	fmt.Fprintln(w, "Total unique browsers", *uniqueBrowsers)
}

func parseUser2(user structs.User, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp, foundUsers *string, i int) {
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
	*foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
}
