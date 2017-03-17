package main

import "sync"

type RoundRobinBalancer struct {
	sync.Mutex
	count   int
	servers []int
	next    int
}

func (b *RoundRobinBalancer) Init(n int) {
	b.count = n
	b.servers = make([]int, n)
	b.next = 0
}

func (b *RoundRobinBalancer) GiveStat() []int {
	b.Lock()
	defer b.Unlock()
	return b.servers
}

func (b *RoundRobinBalancer) GiveNode() int {
	b.Lock()
	defer b.Unlock()
	res := b.next
	b.servers[res]++
	b.next = (b.next + 1) % b.count
	return res
}
