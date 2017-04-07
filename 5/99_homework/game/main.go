package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/telegram-bot-api.v4"
)

var players []*Player
var startPosition *Room
var actions []string
var admins map[int]bool

var db *sql.DB

func getCommandsInfo(player *Player) {
	s := "Возможные команды:\n"
	for _, act := range actions {
		s += act + "\n"
	}
	player.HandleOutput(s)
}

func getHiCommand(player *Player) {
	player.HandleOutput("Вы в игре!")
}

var adminButtons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "/help"},
	tgbotapi.KeyboardButton{Text: "осмотреться"},
	tgbotapi.KeyboardButton{Text: "Сброс игры"},
}

var defaultButtons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "/help"},
	tgbotapi.KeyboardButton{Text: "осмотреться"},
}

func FatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://alexp-go-game.herokuapp.com/bot"

func main() {
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("367628029:AAEgmLPXWDhXvXH6WH9jxwBRU5cbhantK5k")
	FatalOnErr(err)

	bot.Debug = true

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	FatalOnErr(err)

	err = db.Ping()
	FatalOnErr(err)

	rows, err := db.Query(`CREATE TABLE IF NOT EXISTS commands (
	  "id" serial NOT NULL,
	  "from" text NOT NULL,
	  "command" text NOT NULL,
	  "result" text NOT NULL,
	  "time" integer NOT NULL
	);`)
	FatalOnErr(err)
	rows.Close()

	// Отдаем инфу по БД
	go sessionFunc("https://alexp-go-game.herokuapp.com:" + port)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем вебхук
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	FatalOnErr(err)

	initGame()

	updates := bot.ListenForWebhook("/bot")
	go http.ListenAndServe(":"+port, nil)

	// Удаление игрока, при его бездействии
	go func() {
		for {
			time.Sleep(time.Minute)
			for i, pl := range players {
				if time.Now().Sub(pl.lastMessageTime) > time.Minute*15 {
					pl.DeletePlayer()
					players = append(players[:i], players[i+1:]...)
				}
			}
		}
	}()

	// получаем все обновления из канала updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Println("received text: ", update.Message.Text)
		name := update.Message.From.UserName
		id := update.Message.From.ID
		if name == "" {
			name = update.Message.From.FirstName + update.Message.From.LastName
		}
		player, ok := findPlayer(name)

		if !ok {
			player = NewPlayer(id, name, update.Message.Chat.ID)

			// Обработка вывода
			go func(player *Player, id int) {
				var message tgbotapi.MessageConfig
				output := player.GetOutput()
				for msg := range output {
					logMes, err := json.Marshal(&LogJSON{Player: player.name, Message: msg})
					FatalOnErr(err)

					var lastInsertID int
					err = db.QueryRow(`INSERT INTO commands ("from", "command", "result", "time") VALUES ($1,$2,$3,$4) RETURNING id;`,
						player.name, player.curCommand, string(logMes), time.Now().Unix()).Scan(&lastInsertID)
					FatalOnErr(err)

					message = tgbotapi.NewMessage(player.chatID, msg)
					if admins[id] {
						message.ReplyMarkup = tgbotapi.NewReplyKeyboard(adminButtons)
					} else {
						message.ReplyMarkup = tgbotapi.NewReplyKeyboard(defaultButtons)
					}
					bot.Send(message)
				}

			}(player, id)

			addPlayer(player)
		}

		player.lastMessageTime = time.Now()

		// Обработка ввода
		go func(update tgbotapi.Update) {
			player.curCommand = update.Message.Text
			switch update.Message.Text {
			case "/start":
				getHiCommand(player)
				fallthrough
			case "/help":
				getCommandsInfo(player)
			case "Сброс игры":
				if admins[id] {
					for _, pl := range players {
						pl.DeletePlayer()
					}
					initGame()
				} else {
					getCommandsInfo(player)
				}
			default:
				player.HandleInput(update.Message.Text)
			}
		}(update)
	}
}
