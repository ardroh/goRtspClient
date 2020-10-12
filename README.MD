# goRtspClient

goRtspClient is golang library for dealing with the RTSP communication.

## Installation

```bash
go get https://github.com/ardroh/goRtspClient
```

## Usage

```go
rtspClient := goRtspClient.InitRtspClient("192.168.0.10", 554, "some-path")

rtspClient.connect()
//(...)
rtspClient.disconnect()
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)