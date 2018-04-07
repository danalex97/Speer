package interfaces

type DHTQuery struct {
  key   string // the key of the node
  size  int    // size of key to be transfered in MB
  node  string // the node which sends/stores the query
  store bool   // store/retrieve
}
