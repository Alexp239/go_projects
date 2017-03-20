package main

var players []*Player
var positions PositionsStruct

func NewPlayer(name string) *Player {
	player := Player{
		name:     name,
		messages: make(chan string),
	}
	players = append(players, &player)
	return &player
}

func addPlayer(player *Player) {
	player.initRooms()
	player.initTasks()
}

func initGuestsRooms() {
	positions.positions = PositionsMap{}
	positions.positions["кухня"] = make(map[string]*Player)
	positions.positions["коридор"] = make(map[string]*Player)
	positions.positions["комната"] = make(map[string]*Player)
	positions.positions["улица"] = make(map[string]*Player)
}

func initGame() {
	initGuestsRooms()
}
