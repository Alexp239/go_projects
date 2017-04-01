/*
В примере 5/21_cache_example
Есть пример с тегированным кешем и корректным перестроением кеша

Надо совместить эти 2 примера в 1 и сделать удобную функцию
func TcacheGet(mkey string, buildCache funcToRebuild() (interface{}, []slice) )

В неё мы передаём ключ mkey и функцию для перестроения buildCache
buildCache возвращает данные которые мы хотим сохранить и список тегов, которые надо заполнить
( и сбрасывать кеш, если они отличаются).
 Если тега такого ещё нету в кеше, то надо создать его со значением 1.
*/

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
)

var c redis.Conn

type CacheItem struct {
	Data interface{} // можем класть любые данные
	Tags map[string]int
}

func getCachedData(mkey string) (string, error) {
	println("redis get", mkey)
	// получает запись, https://redis.io/commands/get
	data, err := c.Do("GET", mkey)
	item, err := redis.String(data, err)

	// если записи нету, то для этого есть специальная ошибка, её надо обрабатывать отдеьно, это почти штатная ситуация, а не что-то страшное
	if err == redis.ErrNil {
		fmt.Println("Record not found in redis (return value is nil)")
		return "", redis.ErrNil
	} else if err != nil {
		return "", err
	}
	return item, nil
}

func buildCacheNow(mkey string, buildCache func() (interface{}, []string)) *CacheItem {
	res := &CacheItem{Tags: map[string]int{}}
	redis.String(c.Do("SET", mkey+"_lock", ""))
	println("lock set")
	data, tags := buildCache()
	res.Data = data
	for _, tag := range tags {
		cur_value, err := redis.Int(c.Do("GET", tag))
		if err == redis.ErrNil {
			c.Do("SET", tag, 1)
			cur_value = 1
		}
		res.Tags[tag] = cur_value
	}

	jsonData, _ := json.Marshal(*res)
	println("json to store: ", string(jsonData))
	redis.String(c.Do("SET", mkey, jsonData))

	n, err := redis.Int(c.Do("DEL", mkey+"_lock"))
	PanicOnErr(err)
	println("lock deleted:", n)
	return res
}

func TcacheGet(mkey string, buildCache func() (interface{}, []string)) *CacheItem {
	var topCache string
	var err error
	cItems := &CacheItem{}

	for j := 0; j < 4; j++ {
		topCache, err = getCachedData(mkey)
		if err == redis.ErrNil {
			// пытаемся сказать "я строю этот кеш, другие - ждите"
			lockStatus, _ := redis.String(c.Do("SET", mkey+"_lock", "", "EX", 3, "NX"))
			if lockStatus != "OK" {
				// кто-то другой держит лок, подождём и попробуем получить запись еще раз
				println("sleep", j)
				time.Sleep(time.Millisecond * 10)
			} else {
				// успешло залочились, можем строить кеш
				break
			}
		} else {
			break
		}
	}

	if err == redis.ErrNil {
		return buildCacheNow(mkey, buildCache)
	} else {
		_ = json.Unmarshal([]byte(topCache), cItems)
		fmt.Printf("top Cache unpacked %+v\n", *cItems)
		keys := make([]interface{}, 0)
		toCompare := make([]int, 0)
		for key, val := range cItems.Tags {
			keys = append(keys, key)
			toCompare = append(toCompare, val)
		}

		// https: //redis.io/commands/mget
		reply, err := redis.Ints(c.Do("MGET", keys...))
		PanicOnErr(err)

		if !reflect.DeepEqual(toCompare, reply) {
			var was_locked, was_find bool
			println("Invalid cache")
			for k := 0; k < 4; k++ {
				lockStatus, _ := redis.String(c.Do("SET", mkey+"_lock", "", "EX", 3, "NX"))
				if lockStatus != "OK" {
					// кто-то другой держит лок, подождём и попробуем получить запись еще раз
					println("sleep", k)
					time.Sleep(time.Millisecond * 10)
					was_locked = true
				} else if err != nil {
					PanicOnErr(err)
				} else {
					was_find = true
					break
				}
			}
			if was_locked && was_find {
				topCache, err = getCachedData(mkey)
				_ = json.Unmarshal([]byte(topCache), cItems)
				return cItems
			} else {
				return buildCacheNow(mkey, buildCache)
			}
		} else {
			return cItems
		}
	}
}

//PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, _ = redis.DialURL("redis://user:@localhost:6379/0")
	defer c.Close()

	mkey := "test"
	c.Do("DEL", mkey)
	c.Do("DEL", "first")
	c.Do("DEL", "second")
	item := &CacheItem{}

	item = TcacheGet(mkey, func() (interface{}, []string) {
		return "res" + mkey, []string{"first", "second"}
	})

	c.Do("INCR", "first")
	fmt.Println("INCR first")

	item = TcacheGet(mkey, func() (interface{}, []string) {
		return "res" + mkey, []string{"first", "second"}
	})

	fmt.Println(item)
}
