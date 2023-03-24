package fanout_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/badu/fanout"
)

func TestNormalUsage(t *testing.T) {
	type Payload struct {
		Name string
	}

	command := fanout.New[Payload]()
	ctx, cancel := context.WithCancel(context.Background())
	g1c, g2c, g3c := false, false, false
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		ch := command.Sub()
		wg.Done()
		for {
			select {
			case <-ch:
				g1c = true
				wg.Done()
			case <-ctx.Done():
				command.Cancel(ch)
				wg.Done()
				return
			}
		}
	}()

	go func() {
		ch := command.Sub()
		wg.Done()
		for {
			select {
			case <-ch:
				g2c = true
				wg.Done()
			case <-ctx.Done():
				command.Cancel(ch)
				wg.Done()
				return
			}
		}
	}()

	go func() {
		ch := command.Sub()
		wg.Done()
		for {
			select {
			case <-ch:
				g3c = true
				wg.Done()
			case <-ctx.Done():
				command.Cancel(ch)
				wg.Done()
				return
			}
		}
	}()
	wg.Wait() // goroutines above needs to be 'ready'

	wg.Add(3) // to wait goroutine receive the payload
	command.Pub(Payload{Name: "Fire at will!"})
	wg.Wait()

	command.Close()

	// test closing two time
	command.Close()

	// test publishing after close
	command.Pub(Payload{Name: "should NOT work"})

	wg.Add(3) // to wait all goroutines exit
	cancel()
	wg.Wait() // wait for goroutines to receive exit

	if !g1c || !g2c || !g3c {
		t.Fatal("one goroutine was not called", g1c, g2c, g3c)
	}
}

func TestSomeGiveUp(t *testing.T) {
	type Payload struct {
		Name string
	}

	command := fanout.New[Payload]()
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	g1calls := 0
	wg.Add(3)
	go func() {
		wg.Done()
		ch := command.Sub()
		for {
			select {
			case <-ch:
				g1calls++
				command.Cancel(ch) // this one gives up after receiving first command
				wg.Done()
				return
			case <-ctx.Done():
				t.Fatal("goroutine1 exit should be never called")
				return
			}
		}
	}()

	g2calls := 0
	go func() {
		wg.Done()
		done := make(chan struct{})
		time.AfterFunc(100*time.Millisecond, func() {
			close(done)
		})
		ch := command.Sub()
		for {
			select {
			case <-ch:
				g2calls++
				t.Fatal("goroutine2 should never receive")
				wg.Done()
			case <-ctx.Done():
				t.Fatal("goroutine2 exit should be never called")
				return
			case <-done:
				command.Cancel(ch)
				return
			}
		}
	}()

	g3calls := 0
	go func() {
		wg.Done()
		ch := command.Sub()
		for {
			select {
			case <-ch:
				g3calls++
				wg.Done()
			case <-ctx.Done():
				command.Cancel(ch)
				wg.Done()
				return
			}
		}
	}()
	wg.Wait() // goroutines above needs to be 'ready'

	<-time.After(500 * time.Millisecond) // wait for goroutine #2 to 'expire'

	wg.Add(2)
	command.Pub(Payload{Name: "Fire at will!"})
	wg.Wait()

	wg.Add(1)
	command.Pub(Payload{Name: "Fire at will #2!"})
	wg.Wait()

	wg.Add(1)
	cancel()
	wg.Wait() // wait for goroutines to receive exit

	if g1calls != 1 {
		t.Log("goroutine #1 should be called once")
	}

	if g2calls != 0 {
		t.Log("goroutine #2 should never been called")
	}

	if g3calls != 2 {
		t.Log("goroutine #3 should be called twice")
	}

}
