package main

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
