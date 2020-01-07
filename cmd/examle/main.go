package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fschnko/sps"
)

func main() {
	bus := sps.New()
	defer bus.Close()

	done := make(chan struct{})

	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go publisher(bus, "top1", "pub1", time.Second, done)
	go publisher(bus, "top1", "pub2", time.Second, done)
	go subscriber(bus, "top1", "sub1", time.Second*2, done)
	go subscriber(bus, "top1", "sub2", time.Second*4, done)

	go publisher(bus, "top2", "pub3", time.Second, done)
	go publisher(bus, "top1", "pub4", time.Second, done)
	go subscriber(bus, "top2", "sub1", time.Second*3, done)
	go subscriber(bus, "top2", "sub2", time.Second*5, done)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	close(done)

	time.Sleep(time.Second)
}

func publisher(bus *sps.Databus, topic, pub string, interval time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	for i := 0; ; i++ {
		select {
		case _, ok := <-done:
			if !ok {
				ticker.Stop()
				return
			}
		case <-ticker.C:
			msg := fmt.Sprintf("%s/#%d", pub, i)
			bus.Publish(topic, []byte(msg))
		}
	}
}

func subscriber(bus *sps.Databus, topic, sub string, interval time.Duration, done <-chan struct{}) {
	bus.Subscribe(topic, sub)
	ticker := time.NewTicker(interval)
	for {
		select {
		case _, ok := <-done:
			if !ok {
				ticker.Stop()
				return
			}
		case <-ticker.C:
			data, err := bus.Poll(topic, sub)
			if err != nil {
				log.Printf("Subscriber %s: poll %s: %v", sub, topic, err)
				continue
			}
			for _, msg := range data {
				fmt.Printf("%s/%s:\t%s\n", sub, topic, string(msg))
			}
		}
	}
}
