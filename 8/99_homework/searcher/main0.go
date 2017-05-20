package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strings"
)

func SearchFile0(filePath string, w http.ResponseWriter, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp) {
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

	users := make([]map[string]interface{}, 0)
	for _, line := range lines {
		user := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			// fmt.Println("cant unmarshal json: ", f.Name(), line, err)
			continue
		}
		users = append(users, user)
	}

	for i, user := range users {
		parseUser0(user, seenBrowsers, uniqueBrowsers, r, &foundUsers, i)
	}

	fmt.Fprintln(w, "found users:\n"+foundUsers)
	fmt.Fprintln(w, "Total unique browsers", *uniqueBrowsers)
}

func parseUser0(user map[string]interface{}, seenBrowsers *[]string, uniqueBrowsers *int, r *regexp.Regexp, foundUsers *string, i int) {
	isAndroid := false
	isMSIE := false

	browsers, ok := user["browsers"].([]interface{})
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

		if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
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
		if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {
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
	email := r.ReplaceAllString(user["email"].(string), " [at] ")
	*foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
}
