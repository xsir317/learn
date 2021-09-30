package main

import (
	"fmt"
	"math/rand"
	"time"

	rc "./rollingcount"
)

func main() {
	var bkt = rc.NewBucket()
	go func() {
		for {
			time.Sleep(time.Millisecond * 88)
			bkt.AddSuccess(rand.Int() % 100)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond * 123)
			bkt.AddFail(rand.Int() % 15)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond * 1000)
			fmt.Println(bkt.RecentSuccessCnt(), bkt.RecentFailCnt(), bkt.RecentSuccRate())
		}
	}()

	time.Sleep(time.Millisecond * 20000)
	bkt.Shutdown()
}
