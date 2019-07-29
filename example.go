package main

import (
	"errors"
	"fmt"
	"github.com/abusizhishen/doOnceWhileConcurrent/src"
	"sync"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		//并发do something
		go doSomeThing()
	}

	time.Sleep(time.Second*5)
}

var once src.JustOnceSameTime

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

	var requestTag = fmt.Sprintf(keyUser, userId)
	if !once.Req(requestTag) {
		fmt.Println("没抢到锁，等待抢到锁的线程执行结束。。。")
		once.Wait(requestTag)
		fmt.Println("等待结束:", keyUser)
		return userCache.GetUser(userId)
	}

	//得到资源后释放锁
	defer once.Release(requestTag)
	fmt.Println(requestTag, "获得锁，let's Go")

	//为演示效果，sleep
	time.Sleep(time.Second * 3)

	//redis读取用户信息
	fmt.Println("redis读取用户信息:", userId)
	user, err = getUserInfoFromRedis(userId)
	if err != nil {
		return
	}

	//用户写入缓存
	fmt.Println("用户写入缓存:",userId)
	userCache.setUser(user)
	return
}

//用户信息缓存
type UserCache struct {
	Users map[int]UserInfo
	sync.RWMutex
}

type UserInfo struct {
	Id     int
	Name   string
	Age    int
	Gender int
	Img    string
}

var userCache UserCache
var errUserNotFound = errors.New("user not found in cache")

func (c *UserCache) GetUser(id int) (user UserInfo, err error) {
	c.RLock()
	var ok bool
	user, ok = userCache.Users[id]
	if ok {
		return
	}

	c.RUnlock()
	return user, errUserNotFound
}

func (c *UserCache) setUser(user UserInfo) {
	c.Lock()
	if c.Users == nil {
		c.Users = make(map[int]UserInfo)
	}

	c.Users[user.Id] = user
	c.Unlock()
	return
}

func getUserInfoFromRedis(id int) (user UserInfo, err error) {
	// ...

	user = UserInfo{
		Id:   12345,
		Name: "someone",
		Age:  18,
	}
	return
}
