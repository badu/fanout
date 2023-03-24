# Fan Out

Fan Out design pattern, using channels with generics.

## What problem does it solve?

Fan Out pattern is a software architecture pattern that involves distributing data (messages / payloads) from a single
sender (publisher) to multiple receivers (listeners).

Each recipient receives a copy of the message - if your `payload` is a pointer, that would violate the pattern. 

The Fan Out pattern is useful in situations where multiple components of a system need to `react` to the same `event` or
receive the same `data`.

This is a naive implementation of the Fan Out pattern, using channels and generics. 

In its shortest form, it looks like this:

```go
package mypack

func FanOut[T any](from <-chan T, to ...chan<- T) {
	for v := range from {
		for _, ch := range to {
			ch <- v
		}
	}
}

```

## Usage

`go get github.com/badu/fanout`

T. B. D.
