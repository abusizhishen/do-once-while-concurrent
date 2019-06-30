package justOnceWhileCocurrent

import (
	"fmt"
	"time"
)

func main() {
	for i:=0;i<10;i++{
		go add(i)
	}

	for i:=0;i<20;i++{
		fmt.Println(n)
		time.Sleep(time.Second)
	}
}

var n int
var once JustOnceSameTime


//example for usage
// i 是
func add(requestIdentifie interface{})  {
	if once.Req(requestIdentifie){
		fmt.Println(requestIdentifie,"获得锁，Go")
		//为演示效果，sleep15秒
		time.Sleep(time.Second*15)
		n++

		//得到资源后释放锁
		once.Release(requestIdentifie)
	}else {
		fmt.Println(requestIdentifie,"没抢到锁，阻塞中。。。")
		once.Wait(requestIdentifie)
		fmt.Println(requestIdentifie,"等待结束，获取结果为:",n)
	}
}

