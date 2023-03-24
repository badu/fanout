package fanout_test

import (
	"context"
	"sync"
	"testing"

	"github.com/badu/fanout"
)

type Payload struct {
	Name string
}

type Publisher struct {
	bus *fanout.Fanner[Payload]
}

func NewPublisher() Publisher {
	return Publisher{bus: fanout.New[Payload]()}
}

func (p *Publisher) Publish() {
	p.bus.Pub(Payload{Name: "Hello"})
}

func (p *Publisher) Bus() *fanout.Fanner[Payload] {
	return p.bus
}

type Consumer struct {
	wg sync.WaitGroup
}

func NewConsumer(ctx context.Context, bus *fanout.Fanner[Payload], wg sync.WaitGroup) Consumer {
	result := Consumer{wg: wg}
	go result.StartListen(ctx, bus)
	return result
}

func (c *Consumer) StartListen(ctx context.Context, bus *fanout.Fanner[Payload]) {
	ch := bus.Sub()
	c.wg.Add(1)

	for {
		select {
		case payload := <-ch:
			if payload.Name != "Hello" {
				panic("supposed to say `Hello` but said " + payload.Name)
			}
		case <-ctx.Done():
			c.wg.Done()
			return
		}
	}
}

func BenchmarkOneToOne(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	pub := NewPublisher()
	NewConsumer(ctx, pub.Bus(), sync.WaitGroup{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish()
	}
	cancel()
}

func BenchmarkOneToTen(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	pub := NewPublisher()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		NewConsumer(ctx, pub.Bus(), wg)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish()
	}
	cancel()
	wg.Wait()
}

func BenchmarkOneToOneHundred(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	pub := NewPublisher()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		NewConsumer(ctx, pub.Bus(), wg)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pub.Publish()
	}
	cancel()
	wg.Wait()
}
