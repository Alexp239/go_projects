package main

import "strings"

var player Player
var actions map[string]func(...string) string
var tasks []Task

// Task - задача
type Task struct {
	name   string
	action func() bool
}

// ItemSlice - набор предметов
type ItemSlice []*Item

// InventoryMap - набор предметов в рюказке
type InventoryMap map[string]*Item

// FurnitureSlice - набор мебели
type FurnitureSlice []*Furniture

// PathSlice - набор путей
type PathSlice []*Path

// Item - предмет
type Item struct {
	name string
}

// Furniture - мебель
type Furniture struct {
	name        string
	description string
	items       ItemSlice
}

// AddItem - добавить предмет в мебель
func (f *Furniture) AddItem(i *Item) {
	f.items = append(f.items, i)
}

// TakeItem - положить предмет в рюкзак
func (f *Furniture) TakeItem(pos int) {
	player.inventory[f.items[pos].name] = f.items[pos]
	f.items = append(f.items[:pos], f.items[pos+1:]...)
}

// DressBag - Одеть рюкзак
func (f *Furniture) DressBag(pos int) {
	player.inventory = InventoryMap{}
	f.items = append(f.items[:pos], f.items[pos+1:]...)
}

// Locker - "замок"
type Locker struct {
	Furniture
	locked     bool
	unlockItem *Item
	openDiscr  string
	closeDiscr string
}

// Unlock - открывает замок
func (l *Locker) Unlock() {
	l.locked = false
}

/* // Lock - закрывает замок
func (l *Locker) Lock() {
	l.locked = true
}*/

// Room - комната
type Room struct {
	name      string
	lookDiscr string
	goDiscr   string
	freeDiscr string
	furniture FurnitureSlice
	paths     PathSlice
}

// AddPath - добавление прохода из одной комнаты в другую
func (r *Room) AddPath(rNew *Room) {
	p1 := Path{room: rNew}
	p2 := Path{room: r}
	r.paths = append(r.paths, &p1)
	rNew.paths = append(rNew.paths, &p2)
}

// AddPathLocked - добавление прохода из одной комнаты в другую
func (r *Room) AddPathLocked(rNew *Room, locker *Locker) {
	p1 := Path{
		room:   rNew,
		locker: locker,
	}
	p2 := Path{
		room:   r,
		locker: locker,
	}
	r.paths = append(r.paths, &p1)
	rNew.paths = append(rNew.paths, &p2)
}

// AddFurniture - добавить предмет в мебель
func (r *Room) AddFurniture(f *Furniture) {
	r.furniture = append(r.furniture, f)
}

// Path - путь из одного места в другое
type Path struct {
	room   *Room
	locker *Locker
}

// Player - игрок
type Player struct {
	name      string
	inventory InventoryMap
	position  *Room
}

// AddNeighbRooms - информация о том, куда можно пройти
// Нужно что-то сделать с улицей
func AddNeighbRooms(r *Room) string {
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

// ShowLockers - информация о "замках"
func ShowLockers(r *Room) (string, bool) {
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

func checkTasks() string {
	res := ""
	fl := false
	for _, task := range tasks {
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

// Look - осмотреться
func Look(ar ...string) string {
	res := player.position.lookDiscr

	showRes, fl := ShowLockers(player.position)
	res += showRes

	if len(player.position.furniture) > 0 {
		for _, furn := range player.position.furniture {
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
		tasksRes := checkTasks()
		if tasksRes != "" {
			res += ", надо " + tasksRes
		}
	}

	res += "."
	res += AddNeighbRooms(player.position)
	return res
}

// GoRoom - идти
func GoRoom(par ...string) string {
	for _, path := range player.position.paths {
		if path.room.name == par[0] {
			if path.locker != nil && path.locker.locked == true {
				return path.locker.name + path.locker.closeDiscr
			}
			player.position = path.room
			res := path.room.goDiscr
			res += AddNeighbRooms(path.room)
			return res
		}
	}
	return "нет пути в " + par[0]
}

// Take - взять
func Take(par ...string) string {
	if player.inventory == nil {
		return "некуда класть"
	}
	takeItem := par[0]
	for _, furn := range player.position.furniture {
		for i, item := range furn.items {
			if item.name == takeItem {
				furn.TakeItem(i)
				return "предмет добавлен в инвентарь: " + takeItem
			}
		}
	}
	return "нет такого"
}

// Dress - одеть
func Dress(par ...string) string {
	takeItem := par[0]
	if takeItem == "рюкзак" {
		for _, furn := range player.position.furniture {
			for i, item := range furn.items {
				if item.name == takeItem {
					furn.DressBag(i)
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
func Apply(par ...string) string {
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

func handleCommand(s string) string {
	command := strings.Split(s, " ")
	if f, ok := actions[command[0]]; ok {
		return f(command[1:]...)
	}
	return "неизвестная команда"
}

func initActions() {
	actions = map[string]func(...string) string{}
	actions["осмотреться"] = Look
	actions["идти"] = GoRoom
	actions["взять"] = Take
	actions["одеть"] = Dress
	actions["применить"] = Apply
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
	}
	corridor := Room{
		name:    "коридор",
		goDiscr: "ничего интересного.",
	}
	room := Room{
		name:      "комната",
		goDiscr:   "ты в своей комнате.",
		freeDiscr: "пустая комната",
	}
	street := Room{
		name:    "улица",
		goDiscr: "на улице весна.",
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

	player = Player{
		position: &kitchen,
	}
}

func initTasks() {
	tasks = []Task{}

	tasks = append(tasks, Task{"собрать рюкзак", func() bool {
		_, ok := player.inventory["конспекты"]
		return ok
	}})

	tasks = append(tasks, Task{"идти в универ", func() bool {
		return false
	}})
}

func initGame() {
	initRooms()
	initActions()
	initTasks()
}
