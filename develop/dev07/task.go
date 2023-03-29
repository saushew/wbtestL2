package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

var or func(channels ...<-chan interface{}) <-chan interface{}

func main() {
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		done := make(chan interface{})

		numberOfChannels := len(channels)
		if numberOfChannels == 0 {
			return nil
		} else if numberOfChannels == 1 {
			return channels[0]
		}

		go func() {
			if numberOfChannels == 2 {
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			} else {
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-or(append(channels[2:], done)...):
				}
			}
			close(done)
		}()

		return done
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v\n", time.Since(start))
	log.Println(runtime.NumGoroutine())
}
