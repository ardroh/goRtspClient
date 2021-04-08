# goRtspClient

goRtspClient is golang library for dealing with the RTSP communication.

## Version

**Actual version**: 0.1-alpha

## Installation

```TBD```

## Usage

goRtspClient for now supports only __RtpAvpTcp__ with basic and digest authentication.

```go

    connectParams := client.RtspConnectionParams{
            IP:           "192.168.0.5",
            Port:         554,
            Path:         "rtspPath",
            Transport:    headers.RtpAvpTcp,
            Transmission: headers.Unicast,
            Credentials:  auth.Credentials{
                Username: "admin",
                Password: "password"
            },
        }
	rtspClient := RtspClient{}

	connectError := rtspClient.Connect(connectParams)

	if connectError != nil {
		//(...) Error handling
	}

	go doSomethingWithData(rtspClient.DataPacketChan)

	//(...)

	rtspClient.Disconnect()
```

## To Do
- ~~Basic authentication~~
- ~~Digest authentication~~
- ~~RTP over RTSP (interleaved)~~
- Go module release
- Unicast support
- Multicast support

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)