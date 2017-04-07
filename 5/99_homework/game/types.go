package main

import "sync"

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

type PlayersMap map[int]*Player

// Item - предмет
type Item struct {
	name         string
	defaultPlace *Furniture
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
func (f *Furniture) TakeItem(pos int, player *Player) {
	player.inventory[f.items[pos].name] = f.items[pos]
	f.items = append(f.items[:pos], f.items[pos+1:]...)
}

// DressBag - Одеть рюкзак
func (f *Furniture) DressBag(pos int, player *Player) {
	player.inventory = InventoryMap{}
	f.TakeItem(pos, player)
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
	mu        sync.Mutex
	players   PlayersMap
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
