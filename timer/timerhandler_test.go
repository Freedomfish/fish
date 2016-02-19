package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_TimerHandler(t *testing.T) {
	timer := NewTimerHandler()
	defer timer.Stop()
	{
		i := int32(0)
		timer.AfterFunc(time.Second*2, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 2)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Fatal(i)
		}
	}
	{
		i := int32(0)
		timer.AfterFunc(time.Second*3, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 3)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Fatal(i)
		}
	}
}

func Benchmark_TimerHandler(b *testing.B) {
	timer := NewTimerHandler()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.AfterFunc(time.Second, func() {})
	}
}
