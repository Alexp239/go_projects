package main

var players []*Player
var start_position *Room

func NewPlayer(name string) *Player {
	player := Player{
		name:     name,
		messages: make(chan string),
	}
	players = append(players, &player)
	return &player
}

func addPlayer(player *Player) {
	player.initTasks()
	player.position = start_position
	start_position.players[player.name] = player
}

func initRooms() {
	roomTable := Furniture{
		name:        "стол",
		description: "на столе: ",
	}

	keys := Item{"ключи"}
	roomTable.AddItem(&keys)
	conspects := Item{"конспекты"}
	roomTable.AddItem(&conspects)

	roomChair := Furniture{
		name:        "стул",
		description: "на стуле - ",
	}
	roomChair.AddItem(&Item{"рюкзак"})

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

	start_position = &kitchen
}

func initGame() {
	initRooms()
}
