package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

var players []*Player
var startPosition *Room
var actions []string
var admins map[int]bool

func NewPlayer(playerID int, name string, chatID int64) *Player {
	player := Player{
		id:       playerID,
		name:     name,
		messages: make(chan string),
		chatID:   chatID,
	}
	players = append(players, &player)
	return &player
}

func addPlayer(player *Player) {
	player.initTasks()
	player.position = startPosition
	startPosition.players[player.id] = player
}

func initRooms() {
	roomTable := Furniture{
		name:        "стол",
		description: "на столе: ",
	}

	keys := Item{
		name:         "ключи",
		defaultPlace: &roomTable,
	}
	roomTable.AddItem(&keys)
	conspects := Item{
		name:         "конспекты",
		defaultPlace: &roomTable,
	}
	roomTable.AddItem(&conspects)

	roomChair := Furniture{
		name:        "стул",
		description: "на стуле - ",
	}
	roomChair.AddItem(&Item{
		name:         "рюкзак",
		defaultPlace: &roomChair,
	})

	kitchen := Room{
		name:      "кухня",
		lookDiscr: "ты находишься на кухне, на столе чай",
		goDiscr:   "кухня, ничего интересного.",
		players:   PlayersMap{},
	}
	corridor := Room{
		name:    "коридор",
		goDiscr: "ничего интересного.",
		players: PlayersMap{},
	}
	room := Room{
		name:      "комната",
		goDiscr:   "ты в своей комнате.",
		freeDiscr: "пустая комната",
		players:   PlayersMap{},
	}
	street := Room{
		name:    "улица",
		goDiscr: "на улице весна.",
		players: PlayersMap{},
	}

	kitchen.AddPath(&corridor)
	corridor.AddPath(&room)

	door := Locker{
		Furniture: Furniture{
			name:        "дверь",
			description: "",
		},
		locked:     true,
		unlockItem: &keys,
		openDiscr:  " открыта",
		closeDiscr: " закрыта",
	}

	corridor.AddPathLocked(&street, &door)

	room.AddFurniture(&roomTable)
	room.AddFurniture(&roomChair)

	startPosition = &kitchen
}

func initGame() {
	initRooms()
	actions = []string{"/start", "/help", "осмотреться", "идти [комната]",
		"взять [предмет]", "одеть [предмет]", "применить [предмет] [замок]",
		"сказать_игроку [игрок] [сообщение]", "сказать [сообщение]"}
	players = []*Player{}
	admins = map[int]bool{47394442: true}
}

func findPlayer(name string) (player *Player, ok bool) {
	for _, player = range players {
		if player.name == name {
			ok = true
			return
		}
	}
	return
}

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

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://alexp-go-game.herokuapp.com/bot"

func main() {
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("367628029:AAEgmLPXWDhXvXH6WH9jxwBRU5cbhantK5k")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем вебхук
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatal(err)
	}

	initGame()

	updates := bot.ListenForWebhook("/bot")
	go http.ListenAndServe(":"+port, nil)

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

			go func(player *Player, id int) {
				var message tgbotapi.MessageConfig
				output := player.GetOutput()
				for msg := range output {
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

		go func(update tgbotapi.Update) {
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
