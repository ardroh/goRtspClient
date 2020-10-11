package main

import (
	"log"
	"time"
)

func main() {
	rtspClient := RtspClient{
		rtspPath: "axis-media/media.amp",
		cSeq:     0,
		ip:       "172.29.7.51",
		port:     554,
	}
	rtspClient.connect()
	startT := time.Now()
	for {
		packet := <-rtspClient.readPacket
		t := time.Now()
		elapsed := t.Sub(startT)
		if elapsed > 10*time.Second {
			break
		}
		log.Printf("Reading %d bytes in %fs.\n", packet.size, elapsed.Seconds())
	}
	rtspClient.disconnect()
}
