package main

// Look - осмотреться
func (player *Player) Look(ar ...string) string {
	r := player.position
	r.mu.Lock()
	defer r.mu.Unlock()
	res := r.lookDiscr

	showRes, fl := player.ShowLockers()
	res += showRes

	if len(r.furniture) > 0 {
		for _, furn := range r.furniture {
			if len(furn.items) > 0 {
				if fl {
					res += ", "
				}
				fl = true
				res += furn.description
				i := false
				for _, item := range furn.items {
					if i {
						res += ", "
					}
					i = true
					res += item.name
				}
			}
		}
	}
	if fl == false {
		res += player.position.freeDiscr
	}

	if player.position.name == "кухня" {
		tasksRes := player.checkTasks()
		if tasksRes != "" {
			res += ", надо " + tasksRes
		}
	}

	res += "."
	res += player.AddNeighbRooms()
	res += player.AddPlayers()
	return res
}

// GoRoom - идти
func (player *Player) GoRoom(par ...string) string {
	for _, path := range player.position.paths {
		if path.room.name == par[0] {
			if path.locker != nil && path.locker.locked == true {
				return path.locker.name + path.locker.closeDiscr
			}
			r := player.position
			r.mu.Lock()
			defer r.mu.Unlock()
			path.room.mu.Lock()
			defer path.room.mu.Unlock()
			delete(r.players, player.name)
			player.position = path.room
			player.position.players[player.name] = player

			res := path.room.goDiscr
			res += player.AddNeighbRooms()
			return res
		}
	}
	return "нет пути в " + par[0]
}

// Take - взять
func (player *Player) Take(par ...string) string {
	r := player.position
	r.mu.Lock()
	defer r.mu.Unlock()
	if player.inventory == nil {
		return "некуда класть"
	}
	takeItem := par[0]
	for _, furn := range r.furniture {
		for i, item := range furn.items {
			if item.name == takeItem {
				furn.TakeItem(i, player)
				return "предмет добавлен в инвентарь: " + takeItem
			}
		}
	}
	return "нет такого"
}

// Dress - одеть
func (player *Player) Dress(par ...string) string {
	r := player.position
	r.mu.Lock()
	defer r.mu.Unlock()
	takeItem := par[0]
	if takeItem == "рюкзак" {
		for _, furn := range r.furniture {
			for i, item := range furn.items {
				if item.name == takeItem {
					furn.DressBag(i, player)
					return "вы одели: " + takeItem
				}
			}
		}
	} else {
		return "неизвестная команда"
	}
	return "нет такого"
}

// Apply - применить
func (player *Player) Apply(par ...string) string {
	r := player.position
	r.mu.Lock()
	defer r.mu.Unlock()
	applyItem := par[0]
	applyFurn := par[1]

	var item *Item
	var ok bool

	if item, ok = player.inventory[applyItem]; !ok {
		return "нет предмета в инвентаре - " + applyItem
	}

	for _, path := range player.position.paths {
		if path.locker != nil && path.locker.name == applyFurn {
			if path.locker.unlockItem == item && path.locker.locked == true {
				path.locker.Unlock()
				return path.locker.name + path.locker.openDiscr
			}
		}
	}

	return "не к чему применить"
}

func (player *Player) SayPlayer(par ...string) {
	r := player.position
	r.mu.Lock()
	defer r.mu.Unlock()
	reciever, ok := r.players[par[0]]
	if !ok {
		player.HandleOutput("тут нет такого игрока")
		return
	}
	if len(par) == 1 {
		reciever.HandleOutput(player.name + " выразительно молчит, смотря на вас")
		return
	}
	res := player.name + " говорит вам:"
	for _, word := range par[1:] {
		res += " " + word
	}
	reciever.HandleOutput(res)
	return
}

func (player *Player) Say(par ...string) {
	r := player.position

	r.mu.Lock()
	defer r.mu.Unlock()

	res := ""
	if len(par) == 0 {
		res += player.name + " выразительно молчит"
	} else {
		res += player.name + " говорит:"
		for _, word := range par {
			res += " " + word
		}
	}

	for _, pl := range r.players {
		pl.HandleOutput(res)
	}
	return
}
