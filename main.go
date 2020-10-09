package main

func main() {
	rtspClient := RtspClient{
		rtspPath: "axis-media/media.amp",
		cSeq:     0,
		ip:       "172.29.7.51",
		port:     554,
	}
	rtspClient.connect()
}
