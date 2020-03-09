package main

import (
	"errors"
	"fmt"
	"github.com/abusizhishen/do-once-while-concurrent/src"
	"log"
	"sync"
	"time"
)

func main() {
	//并发do something
	for i := 0; i < 5; i++ {
		go doSomeThing()
	}

	time.Sleep(time.Second * 5)
}

var once src.DoOnce

func doSomeThing() {
	var userId = 12345
	var user, err = getUserInfo(userId)
	fmt.Println(user, err)
}

//example for usage
// 演示获取用户详情的过程，先从本地缓存读取用户,如果本地缓存不存在,就从redis读取
var keyUser = "user_%d"

func getUserInfo(userId int) (user UserInfo, err error) {
	user, err = userCache.GetUser(userId)
	if err == nil {
		return
	}

	log.Println(err)
	var requestTag = fmt.Sprintf(keyUser, userId)
	if !once.Req(requestTag) {
		log.Println("没抢到锁，等待抢到锁的线程执行结束。。。")
		once.Wait(requestTag)
		log.Println("等待结束:", requestTag)
		return userCache.GetUser(userId)
	}

	//得到资源后释放锁
	defer once.Release(requestTag)
	log.Println(requestTag, "获得锁，let's Go")

	//为演示效果，sleep
	time.Sleep(time.Second * 3)

	//redis读取用户信息
	log.Println("redis读取用户信息:", userId)
	user, err = getUserInfoFromRedis(userId)
	if err != nil {
		return
	}

	//用户写入缓存
	log.Println("用户写入缓存:", userId)
	userCache.setUser(user)
	return
}

//用户信息缓存
type UserCache struct {
	Users map[int]UserInfo
	sync.RWMutex
}

type UserInfo struct {
	Id   int
	Name string
	Age  int
}

var userCache UserCache
var errUserNotFound = errors.New("user not found in cache")

func (c *UserCache) GetUser(id int) (user UserInfo, err error) {
	c.RLock()
	defer c.RUnlock()
	var ok bool
	user, ok = userCache.Users[id]
	if ok {
		return
	}

	return user, errUserNotFound
}

func (c *UserCache) setUser(user UserInfo) {
	c.Lock()
	defer c.Unlock()
	if c.Users == nil {
		c.Users = make(map[int]UserInfo)
	}

	c.Users[user.Id] = user
	return
}

func getUserInfoFromRedis(id int) (user UserInfo, err error) {
	user = UserInfo{
		Id:   12345,
		Name: "abusizhishen",
		Age:  18,
	}
	return
}
