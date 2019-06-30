package justOnceWhileCocurrent

import (
	"fmt"
	"github.com/abusizhishen/justOnceWhileCocurrent/src"
	"time"
)

func main() {
	var requestIdentifie = "382h23ehd32je32"
	for i:=0;i<10;i++{
		go add(requestIdentifie)
	}

	for i:=0;i<5;i++{
		fmt.Println(n)
		time.Sleep(time.Second)
	}
}

var n int
var once src.JustOnceSameTime

//example for usage
// i 是
func add(requestIdentifie interface{})  {
	if once.Req(requestIdentifie){
		fmt.Println(requestIdentifie,"获得锁，Go")
		//为演示效果，sleep
		time.Sleep(time.Second*3)
		n++

		//得到资源后释放锁
		once.Release(requestIdentifie)
	}else {
		fmt.Println("没抢到锁，等待抢到锁的线程执行结束。。。")
		once.Wait(requestIdentifie)
		fmt.Println("等待结束，获取结果为:",n)
	}
}