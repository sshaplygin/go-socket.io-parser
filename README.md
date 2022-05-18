# go-socket.io-parser

## Documentation

### Encoding format:
```
<packet type>[<count of binary attachments>-][<namespace>,][<acknowledgment id>][JSON-stringified payload without binary]
+ binary attachment extracted
+ binary attachment extracted
```

packet type consists of header + data.

`socketio_parser.encode()` -> `go_socketio_parser.Marshal(packet *parser.Packet) ([]byte, error)` <br/>
`socketio_parser.decode()` -> `go_socketio_parser.Unmarshal(data []byte, packet *parser.Packet) error` <br/>


### Methods:

same approach as `encoding/json`:

Encode by custom writer:
```go
err := go_socketio_parser.NewEncoder(w io.Writer).Encode(packet *Packet)
```

Decode by custom reader:
```go
err := go_socketio_parser.NewEncoder(r io.Reader).Decode(packet *Packet)
```

## TODO

* Add validate test cases for invalid payload (link)[https://github.com/socketio/socket.io-parser/blob/main/test/parser.js#L134]
* Add streaming API
* Add inner structs
* Unit tests

## Benchmarks

### Compare changes
```bash
go test -run=NONE -bench=. ./... > old.txt
# make changes
go test -run=NONE -bench=. ./... > new.txt

benchcmp old.txt new.txt
```

### Bench with profiles

```bash
go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out ./...
```

## Links

- [parser](https://github.com/socketio/socket.io-parser)
- [docs](https://github.com/socketio/socket.io-protocol)