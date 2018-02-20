package model

type Join struct {
  nodeId string
}

func NewJoin(nodeId string) *Join {
  j := new(Join)
  j.nodeId = nodeId
  return j
}


func (j *Join) NodeId() string {
  return j.nodeId
}
