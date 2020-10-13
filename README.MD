# goRtspClient

goRtspClient is golang library for dealing with the RTSP communication.

## Installation

```bash
go get https://github.com/ardroh/goRtspClient
```

## Usage

goRtspClient for now supports only __RtpAvpTcp__ and __Unicast__.

```go
rtspClient := goRtspClient.InitRtspClient(goRtspClient.RtspConnectionParams{
    IP:           "192.168.0.10",
    Port:         554,
    Path:         "some-path",
    Transport:    commands.RtpAvpTcp,
    Transmission: commands.Unicast,
})

connectError := rtspClient.connect()
if connectError != nil {
    //(...) error handling
}
for {
        packet := <-rtspClient.GetReadChan()
        //(...) do something with a packet
}
//(...)
rtspClient.disconnect()
```

## To Do
- Basic authentication
- Digest authentication
- Unicast support
- Multicast support

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)