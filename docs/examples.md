# Examples

### Examples

Before running the examples, move to the installation directory of Speer:
```
cd $GOPATH/src/github.com/danalex97/Speer
```

We provide the following examples:
```
go run speer.go -config=examples/config/sink.json
go run speer.go -config=examples/config/broadcast.json
go run speer.go -config=examples/config/data_link.json
```

For a more complex example, check this [repository](https://github.com/danalex97/nfsTorrent).

### Primitives

Getting the bootstrap node id:
```go
id := util.Join()
```

Sending a control message:

```go
util.Transport().ControlSend(id, "message")
```

Receiving a control message:

```go
msg := <-util.Transport().ControlRecv()
```

Sending data via a link:

```go
// Creating the link
link := util.Transport().Connect(id)

// Sending the data
link.Upload(Data{
  Id   : "someUniqueId", // Some ID associated with the message.
                         // The ID can be used for keeping the actual data or metadata.
  Size : 1000, // Total data size in bits.
})
```

Getting data from a link:

```go
data := <-link.Download()
```
