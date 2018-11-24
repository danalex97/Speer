package config

type NodeConf struct {
  Number    int `json:"number"`
  Upload    int `json:"upload"`
  Download  int `json:"download"`
}

type Config struct {
  // the language of the entry point and the path of the entry point: file
  // which will be used as main class running in the simulation
  Lang string `json:"lang"`
  Entry string `json:"entry"`

  // number of peers in the system
  Nodes int `json:"nodes"`

  // network level topology parameters
  TransitDomains     int `json:"transitDomains"`
  TransitDomainSize  int `json:"transitDomainSize"`
  StubDomains        int `json:"stubDomains"`
  StubDomainSize     int `json:"stubDomainSize"`

  // interval of virtual time units at which the capacity scheduler runs
  TransferInterval   int `json:"transferInterval"`

  // list which shows the distribution of node capacities
  CapacityNodes      []NodeConf `json:"capacityNodes"`

  // Latency support - when latency support is off, nodes run directly
  Latency bool `json:"latency"`

  // Run the simulator's event queue with support for parallel events.
  Parallel bool `json:"parallel"`
}
