package src

import (
	"sync"
)

type DoOnce struct {
	lock sync.RWMutex
	data data
}

type data map[interface{}]*sync.WaitGroup

func New() *DoOnce {
	return &DoOnce{
		lock: sync.RWMutex{},
		data: make(data),
	}
}

/*
RequestTag 请求标识 用于标识同一个资源
同一时刻只有一个请求能获取执行权限，获得执行权限的线程接下来需要执行具体的业务逻辑，
完成后调用release方法通知其他线程，操作已完成，获取资源即可
其他请求接下来需要调用wait方法
*/
func (u *DoOnce) Req(RequestTag interface{}) bool {
	u.lock.Lock()
	defer u.lock.Unlock()

	if u.data == nil {
		u.data = make(data)
	}

	_, ok := u.data[RequestTag]
	if ok {
		//log.Println("没有得到锁，等待执行者执行结束")
		return false
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	u.data[RequestTag] = wg
	//log.Println("获取锁:", RequestTag)

	return true
}

/*RequestTag 请求标识 用于标识同一个资源
调用wait方法将处于阻塞状态，直到获得执行权限的线程处理完具体的业务逻辑，调用release方法来通知其他线程资源ok了
*/
func (u *DoOnce) Wait(RequestTag interface{}) {
	u.lock.RLock()
	w, ok := u.data[RequestTag]
	u.lock.RUnlock()
	if !ok {
		//log.Println("等待结束：", RequestTag)
		return
	}

	w.Wait()
	//log.Println("等待结束：", RequestTag)
	return
}

/*RequestTag 请求标识 用于标识同一个资源
获得执行权限的线程需要在执行完业务逻辑后调用该方法通知其他处于阻塞状态的线程
*/
func (u *DoOnce) Release(RequestTag interface{}) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.data[RequestTag]; !ok {
		//log.Println("锁已释放？还是不存在？RequestTag用错？RequestTag: ", RequestTag)
		return
	}
	u.data[RequestTag].Done()
	delete(u.data, RequestTag)
	//log.Println("释放锁:", RequestTag)
}
