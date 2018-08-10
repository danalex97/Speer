package bridge

type Message interface {
  GetId() string
}

const (
  CREATE = iota
)

type Create struct {
  Id string `json:"id"`
}

func (c *Create) GetId() string {
  return c.Id
}

func MessageToType(m Message) int {
  switch m.(type) {
  case *Create:
    return CREATE
  }
  panic("Unrecognized message")
}

func TypeToMessage(tp int) Message {
  switch tp {
  case CREATE:
    return &Create{}
  }
  panic("Unrecognized message")
}
