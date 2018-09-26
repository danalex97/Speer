package interfaces

type EnvironBridge interface {
  SendMessage(Message)
  RecvMessage(Id string) <-chan Message
}

type Message interface {
  GetId() string
}
