package main

import "strings"

// Player - игрок
type Player struct {
	name      string
	inventory InventoryMap
	position  *Room
	tasks     []Task
	messages  chan string
}

func (player *Player) HandleInput(s string) {
	command := strings.Split(s, " ")
	switch command[0] {
	case "осмотреться":
		player.HandleOutput(player.Look(command[1:]...))
	case "идти":
		player.HandleOutput(player.GoRoom(command[1:]...))
	case "взять":
		player.HandleOutput(player.Take(command[1:]...))
	case "одеть":
		player.HandleOutput(player.Dress(command[1:]...))
	case "применить":
		player.HandleOutput(player.Apply(command[1:]...))
	case "сказать_игроку":
		player.SayPlayer(command[1:]...)
	case "сказать":
		player.Say(command[1:]...)
	default:
		player.HandleOutput("неизвестная команда")
	}
}

func (player *Player) HandleOutput(s string) {
	player.messages <- s
}

func (player *Player) GetOutput() <-chan string {
	/*c := make(chan string)
	go func() {
		for {
			select {
			case s := <-player.messages:
				c <- s
			}
		}
	}()*/
	return player.messages
}

func (player *Player) checkTasks() string {
	res := ""
	fl := false
	for _, task := range player.tasks {
		if !task.action() {
			if fl {
				res += " и "
			}
			fl = true
			res += task.name
		}
	}
	return res
}

// ShowLockers - информация о "замках"
func (player *Player) ShowLockers() (string, bool) {
	res := ""
	j := false
	for _, path := range player.position.paths {
		if path.locker != nil {
			if j == true {
				res += ", "
			}
			j = true
			res += path.locker.name + " -"
			if path.locker.locked {
				res += path.locker.closeDiscr
			} else {
				res += path.locker.openDiscr
			}
		}
	}
	return res, j
}

// AddNeighbRooms - информация о том, куда можно пройти
// Нужно что-то сделать с улицей
func (player *Player) AddNeighbRooms() string {
	r := player.position
	res := " можно пройти - "
	if r.name == "улица" {
		res += "домой"
		return res
	}
	for i, path := range player.position.paths {
		if i != 0 {
			res += ", "
		}
		res += path.room.name
	}
	return res
}

func (player *Player) AddPlayers() string {
	r := player.position
	res := ""
	if len(r.players) > 1 {
		res += ". Кроме вас тут ещё"
		j := false
		for _, pl := range r.players {
			if pl != player {
				if j {
					res += ","
				}
				j = true
				res += " " + pl.name
			}
		}
	}
	return res
}

func (player *Player) initTasks() {
	player.tasks = []Task{}

	player.tasks = append(player.tasks, Task{"собрать рюкзак", func() bool {
		_, ok := player.inventory["конспекты"]
		return ok
	}})

	player.tasks = append(player.tasks, Task{"идти в универ", func() bool {
		return false
	}})
}
