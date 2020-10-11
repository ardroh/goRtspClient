package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	rtspClient := RtspClient{
		rtspPath: "axis-media/media.amp",
		cSeq:     0,
		ip:       "172.29.7.51",
		port:     554,
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println(sig)
			if sig == os.Interrupt {
				rtspClient.disconnect()
				return
			}
		}
	}()
	rtspClient.connect()
	startT := time.Now()
	for {
		packet := <-rtspClient.readPacket
		t := time.Now()
		elapsed := t.Sub(startT)
		if elapsed > 600*time.Second {
			break
		}
		log.Printf("Reading %d bytes in %fs.\n", packet.size, elapsed.Seconds())
	}
	rtspClient.disconnect()
}
