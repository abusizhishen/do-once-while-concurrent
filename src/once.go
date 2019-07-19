package src

import (
	"log"
	"sync"
)

type JustOnceSameTime struct {
	Lock sync.RWMutex
	Map  map[interface{}]chan bool
}

/*
RequestTag 请求标识 用于标识同一个资源
同一时刻只有一个请求能获取执行权限，获得执行权限的线程接下来需要执行具体的业务逻辑，
完成后调用release方法通知其他线程，操作已完成，获取资源即可
其他请求接下来需要调用wait方法
*/
func (u *JustOnceSameTime) Req(RequestTag interface{}) bool {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	if u.Map == nil {
		u.Map = make(map[interface{}]chan bool)
	}

	_, ok := u.Map[RequestTag]
	if ok {
		log.Println("没有得到锁，等待执行者执行结束")
		return false
	}

	u.Map[RequestTag] = make(chan bool, 1)
	log.Println("获取锁:", RequestTag)

	return true
}

/*RequestTag 请求标识 用于标识同一个资源
调用wait方法将处于阻塞状态，直到获得执行权限的线程处理完具体的业务逻辑，调用release方法来通知其他线程资源ok了
*/
func (u *JustOnceSameTime) Wait(RequestTag interface{}) {
	u.Lock.RLock()
	c, ok := u.Map[RequestTag]
	u.Lock.RUnlock()
	if !ok {
		log.Println("等待结束：", RequestTag)
		return
	}
	select {
	case <-c:
		log.Println("等待结束：", RequestTag)
		return
	}
}

/*RequestTag 请求标识 用于标识同一个资源
获得执行权限的线程需要在执行完业务逻辑后调用该方法通知其他处于阻塞状态的线程
*/
func (u *JustOnceSameTime) Release(RequestTag interface{}) {
	u.Lock.Lock()
	if _, ok := u.Map[RequestTag]; !ok {
		log.Println("锁已释放？还是不存在？RequestTag用错？RequestTag: ", RequestTag)
		u.Lock.Unlock()
		return
	}
	close(u.Map[RequestTag])
	delete(u.Map, RequestTag)
	u.Lock.Unlock()
	log.Println("释放锁:", RequestTag)
}
