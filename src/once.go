package src

import (
	"log"
	"sync"
)

type JustOnceSameTime struct {
	Lock sync.RWMutex
	Map map[interface{}]chan bool
}

func (u *JustOnceSameTime)Req(key interface{}) bool {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	if u.Map == nil{
		u.Map = make(map[interface{}]chan bool)
	}

	_,ok := u.Map[key]
	if ok{
		//log.Println("没有得到锁，等待执行者执行结束")
		return false
	}

	u.Map[key] = make(chan bool,1)
	//log.Println("获取锁")

	return true
}

func (u *JustOnceSameTime)Wait(key interface{}){
	for{
		u.Lock.RLock()
		_,ok := u.Map[key]
		if !ok {
			//log.Println("等待结束：")
			u.Lock.RUnlock()
			return
		}
		select {
		case _,ok:= <- u.Map[key]:
			if !ok {
				//log.Println("等待结束：")
				u.Lock.RUnlock()
				return
			}else{
				log.Println(ok)
			}
		default:
		}
		u.Lock.RUnlock()
	}
}

func (u *JustOnceSameTime)Release(str interface{}){
	u.Lock.Lock()
	close(u.Map[str])
	delete(u.Map,str)
	u.Lock.Unlock()
	//log.Println("释放锁")
}


