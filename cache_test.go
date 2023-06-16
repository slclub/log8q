package log8q

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache(10, 8192)
	wg := sync.WaitGroup{}
	wg_reader := sync.WaitGroup{}
	n, nn, ns := uint64(0), 0, uint64(0)

	readfn := func(ctx context.Context, ww *sync.WaitGroup) {
		dst := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				//fmt.Println("ctx -> done")
				ww.Done()
				return
			default:
				//rn := cache.ReadSize()

				n1, _ := cache.Read(dst)
				nn += n1
			}
		}
	}
	//muw := sync.Mutex{}
	logs := "hello,world! where are you. I am chinese!"
	logs += logs
	logs += logs
	logs += logs
	logs += logs
	logs += logs
	logs += logs
	logs += logs
	logs += logs
	//logs += logs
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		atomic.AddUint64(&ns, 4*uint64(len(logs)))
		go func() {
			for i := 0; i < 4; i++ {
				n1, _ := cache.Write([]byte(logs))
				atomic.AddUint64(&n, uint64(n1))
			}
			//fmt.Print("W:", n1)
			wg.Done()
		}()

	}
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 1; i++ {
		wg_reader.Add(1)
		go readfn(ctx, &wg_reader)
	}
	// write
	wg.Wait()
	cache.ReadSize()

	//time.Sleep(time.Millisecond * 1000)
	// read
	cancel()
	wg_reader.Wait()
	fmt.Println("string lenght is ", ns, "write leght:", n, " read lenght", nn)
}
