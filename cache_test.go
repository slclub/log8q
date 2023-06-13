package log8q

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache(5, 8292)
	wg := sync.WaitGroup{}
	wg_reader := sync.WaitGroup{}
	n, nn, ns := 0, 0, 0

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
	logs := "hello,world! where are you. I am chinese"
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
		ns += len(logs)
		go func() {
			//muw.Lock()
			n1, _ := cache.Write([]byte(logs))
			//muw.Unlock()
			n += n1
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
