# go-socket.io-parser

## Documentation

### Encoding format:
```
<packet type>[<# of binary attachments>-][<namespace>,][<acknowledgment id>][JSON-stringified payload without binary]

+ binary attachments extracted.
```

`socketio_parser.encode()` -> `go_socketio_parser.Marshal(h Header, attach []interface{}) ([]byte, error)` <br/>
`socketio_parser.decode()` -> `go_socketio_parser.Unmarshal(data []byte, h Header, attach []interface{}) error` <br/>


### Methods:
same approach as `encoding/json`:

Encode by custom writer:
```go
err := go_socketio_parser.NewEncoder(w io.Writer).Encode(h Header, attach interface{})
```

Decode by custom reader:
```go
err := go_socketio_parser.NewEncoder(r io.Reader).Decode(h Header, attach interface{})
```

## TODO

* Refactoring code
* Refactoring tests
* Add validate test cases for invalid payload (link)[https://github.com/socketio/socket.io-parser/blob/main/test/parser.js#L134]

## Links

- [parser](https://github.com/socketio/socket.io-parser)
- [docs](https://github.com/socketio/socket.io-protocol)