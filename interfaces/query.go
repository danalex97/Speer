package interfaces

type Query interface {
  Key()   string // the key of the node
  Size()  int    // size of key to be transfered in MB
  Node()  string // the node which sends/stores the query
  Store() bool   // store/retrieve
}
