# Fan Out

Fan Out design pattern, using channels with generics.

## What problem does it solve?

Fan Out pattern is a software architecture pattern that involves distributing data (messages / payloads) from a single
sender (publisher) to multiple receivers (listeners).

Each recipient receives a copy of the message - if your `payload` is a pointer, that would violate the pattern. 

The Fan Out pattern is useful in situations where multiple components of a system need to `react` to the same `event` or
receive the same `data`.

## Usage

`go get github.com/badu/fanout`

T.B.D.
