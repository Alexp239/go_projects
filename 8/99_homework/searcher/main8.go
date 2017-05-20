package main

import (
	// "encoding/json"
	"bufio"
	"fmt"
	jlexer "github.com/mailru/easyjson/jlexer"
	"msu-go-11/8/99_homework/searcher/structs"
	_ "net/http/pprof"
	"os"
	"strings"
)

func SearchFile8(fileInf *FileInfo) {
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
		parseUser8(user, fileInf, i)
		i++
	}

	fmt.Fprintln(fileInf.w, "Total unique browsers", fileInf.uniqueBrowsers)
}

func parseUser8(user structs.User1, fileInf *FileInfo, i int) {
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
	email := strings.Replace(user.Email, "@", " [at] ", -1)
	fmt.Fprintln(fileInf.w, "[%d] %s <%s>\n", i, user.Name, email)
}
