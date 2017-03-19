package main

var players []*Player
var postions PositionsStruct

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
	postions.positions = PositionsMap{}
	postions.positions["кухня"] = make(map[string]*Player)
	postions.positions["коридор"] = make(map[string]*Player)
	postions.positions["комната"] = make(map[string]*Player)
	postions.positions["улица"] = make(map[string]*Player)
}

func initGame() {
	initGuestsRooms()
}
