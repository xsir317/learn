package rollingcount

import (
	"sync"
	"time"
)

//数据结构，一组桶， 我们用2个切片来表示
//size 默认10好了， 不失一般性
//对外提供 RecentSuccessCnt 、 RecentFailCnt 和 RecentSuccRate ， 提供 AddSuccess 和AddFail 就好了。 使用的时候就是
//新建一个对象， 不停告诉它 成功失败成功失败
//想问数量就去问一下
type bucket_counter struct {
	succ_buckets []int
	fail_buckets []int
	//这里锁比较重， 如果针对每一个桶配一个锁好不好呢？
	mutex *sync.RWMutex
	done  chan interface{}
}

func getBucket(now time.Time) int {
	return (int)(now.Unix() % 11)
}

func NewBucket() *bucket_counter {
	var ret_bucket = new(bucket_counter)
	ret_bucket.succ_buckets = make([]int, 11)
	ret_bucket.fail_buckets = make([]int, 11)
	ret_bucket.mutex = new(sync.RWMutex)
	ret_bucket.done = make(chan interface{})

	go func() {
		var next_bucket int
		for {
			time.Sleep(time.Millisecond * 333)
			select {
			case <-ret_bucket.done:
				return
			default:
				next_bucket = getBucket(time.Now().Add(time.Second))
				ret_bucket.mutex.Lock()
				ret_bucket.succ_buckets[next_bucket] = 0
				ret_bucket.fail_buckets[next_bucket] = 0
				ret_bucket.mutex.Unlock()
			}
		}
	}()
	return ret_bucket
}

func (r *bucket_counter) Shutdown() {
	close(r.done)
}

func (r *bucket_counter) AddSuccess(i int) {
	if i == 0 {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.succ_buckets[getBucket(time.Now())] += i
}

func (r *bucket_counter) AddFail(i int) {
	if i == 0 {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.fail_buckets[getBucket(time.Now())] += i
}

func (r *bucket_counter) RecentSuccessCnt() int {

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sum := 0
	next_bucket := getBucket(time.Now().Add(time.Second))
	for _, val := range r.succ_buckets {
		//我们在这里把所有值都加上， 而不是判断下标是否等于 next_bucket ，最后再减掉 next_bucket 。
		//算是一点点优化？
		sum += val
	}
	return sum - r.succ_buckets[next_bucket]
}

func (r *bucket_counter) RecentFailCnt() int {

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sum := 0
	next_bucket := getBucket(time.Now().Add(time.Second))
	for _, val := range r.fail_buckets {
		sum += val
	}
	return sum - r.fail_buckets[next_bucket]
}

func (r *bucket_counter) RecentSuccRate() float64 {
	s := r.RecentSuccessCnt()
	f := r.RecentFailCnt()

	total := s + f
	if total == 0 {
		return 0.0
	}

	return float64(s) / float64(total)
}
