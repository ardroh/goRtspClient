# goRtspClient

goRtspClient is golang library for dealing with the RTSP communication.

## Installation

```bash
go get https://github.com/ardroh/goRtspClient
```

## Usage

```go
    rtspClient := RtspClient{
		rtspPath: "some-path",
		cSeq:     0,
		ip:       "192.168.0.10",
		port:     554,
    }

	rtspClient.connect()
    
    //(...)
    
    rtspClient.disconnect()
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)