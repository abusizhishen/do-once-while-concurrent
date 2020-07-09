package src

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	once := New()
	assert.IsType(t, &DoOnce{}, once)
}

func TestDoOnce_Req(t *testing.T) {
	t.Run("reqIntTag", func(t *testing.T) {
		once := New()
		result := once.Req(123)
		assert.True(t, result)
	})

	t.Run("reqStringTag", func(t *testing.T) {
		once := New()
		result := once.Req("foo")
		assert.True(t, result)
	})

	t.Run("reqSameTagAgain", func(t *testing.T) {
		once := New()
		once.Req(123)
		result := once.Req(123)

		assert.False(t, result)
	})

	t.Run("reqDifferentTag", func(t *testing.T) {
		once := New()
		once.Req(123)
		result := once.Req("foo")

		assert.True(t, result)
	})

	t.Run("reqSameTagAfterRelease", func(t *testing.T) {
		once := New()
		once.Req(123)
		once.Release(123)
		result := once.Req(123)

		assert.True(t, result)
	})

	t.Run("reqAfterReleaseDifferentTag", func(t *testing.T) {
		once := New()
		once.Req(123)
		once.Release(123)
		result := once.Req("foo")

		assert.True(t, result)
	})

	t.Run("performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkDoOnce_Req(1))
		res4 := testing.Benchmark(benchmarkDoOnce_Req(4))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
	})
}

func TestDoOnce_Wait(t *testing.T) {
	t.Run("sameTagWaitWithNotRelease", func(t *testing.T) {
		var once = New()
		var array = new([]interface{})
		var key = "foo"
		once.Req(key)

		go func() {
			once.Wait(key)
			*array = append(*array, 1)
		}()

		time.Sleep(time.Millisecond)
		assert.Equal(t, 0, len(*array))
	})

	t.Run("sameTagWaitWithRelease", func(t *testing.T) {
		var once = New()
		var array = new([]interface{})
		var key = "foo"
		once.Req(key)

		go func() {
			once.Wait(key)
			*array = append(*array, 1)
		}()

		once.Release(key)
		time.Sleep(time.Millisecond)
		assert.Equal(t, 1, len(*array))
	})

	t.Run("differentTagWait", func(t *testing.T) {
		var once = New()
		var array = new([]interface{})
		var key = "foo"
		once.Req(key)

		go func() {
			once.Wait("boo")
			*array = append(*array, 1)
		}()

		time.Sleep(time.Millisecond)
		assert.Equal(t, 1, len(*array))
	})
}

func TestDoOnce_Release(t *testing.T) {
	t.Run("release", func(t *testing.T) {
		var score, result = new(uint32), new(uint32)
		var lock sync.RWMutex

		getScore := func() (uint32, bool) {
			lock.RLock()
			defer lock.RUnlock()
			if *score == 0 {
				return 0, false
			} else {
				return *score, true
			}
		}

		setScore := func() {
			lock.Lock()
			defer lock.Unlock()
			*score = 100
		}

		add := func(result *uint32, score uint32) {
			atomic.AddUint32(result, score)
		}

		key := 123
		once := New()
		var wg sync.WaitGroup
		var count = 100
		wg.Add(count)

		for i := 0; i < count; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				score, exists := getScore()
				if exists {
					add(result, score)
					return
				}

				if !once.Req(key) {
					once.Wait(key)
					score, _ := getScore()
					add(result, score)
				}

				setScore()
				score, _ = getScore()
				add(result, score)
				once.Release(key)
			}(&wg)
		}

		wg.Wait()
		assert.Equal(t, *result, uint32(count)*(*score))
	})
}

func BenchmarkDoOnce_Req(b *testing.B) {
	benchmarkDoOnce_Req(1)(b)
	benchmarkDoOnce_Req_same_key(1)(b)
}

func benchmarkDoOnce_Req(n int) func(*testing.B) {
	once := New()
	return func(b *testing.B) {
		for i := 0; i < 1000*n; i++ {
			once.Req(i)
		}
	}
}

func benchmarkDoOnce_Req_same_key(n int) func(*testing.B) {
	once := New()
	return func(b *testing.B) {
		for i := 0; i < 1000*n; i++ {
			once.Req(1)
		}
	}
}

func BenchmarkAll(b *testing.B) {
	b.Run("benchmarkReq", BenchmarkDoOnce_Req)
}
