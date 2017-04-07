package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type LogJSON struct {
	Player  string `json:"player"`
	Message string `json:"message"`
}

var loginFormTmpl = `
<html>
	<body>
	<form action="/get_cookie" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`

type sessionData struct {
	userName  string
	userAgent string
	lastTime  time.Time
}

var sessions = map[string]sessionData{}

func sessionFunc(addr string) {

	go func() {
		for {
			time.Sleep(time.Minute)
			for sessionID, session := range sessions {
				if time.Now().Sub(session.lastTime) > time.Minute*15 {
					delete(sessions, sessionID)
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		sessionID, err := r.Cookie("session_id")

		if err == http.ErrNoCookie {
			w.Write([]byte(loginFormTmpl))
			return
		} else if err != nil {
			FatalOnErr(err)
		}

		data, ok := sessions[sessionID.Value]
		username := data.userName
		userAgent := data.userAgent

		if !ok {
			sessionID.Expires = time.Now()
			sessionID.MaxAge = -1
			http.SetCookie(w, sessionID)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			curSessionData := sessions[sessionID.Value]
			curSessionData.lastTime = time.Now()
			sessions[sessionID.Value] = curSessionData

			fmt.Fprintln(w, "Welcome, "+username, " from ", userAgent)
			var id, tm int
			var from, com, res string
			rows, err := db.Query("SELECT * FROM commands ORDER BY time DESC LIMIT 100;")
			FatalOnErr(err)
			for rows.Next() {
				err = rows.Scan(&id, &from, &com, &res, &tm)
				FatalOnErr(err)
				fmt.Fprintln(w, id, from, com, res, tm)
			}
			rows.Close()
		}
	})

	http.HandleFunc("/get_cookie", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		passw := "111"

		inputLogin := r.Form["login"][0]
		inputPassword := r.Form["password"][0]

		if inputPassword != passw {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		expiration := time.Now().Add(15 * time.Minute)

		sessionID := RandStringRunes(32)
		userAgent, _ := r.Header["User-Agent"]
		sessions[sessionID] = sessionData{inputLogin, userAgent[0], time.Now()}

		cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(addr, nil)
}
