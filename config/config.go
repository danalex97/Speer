package config

type NodeConf struct {
  Number    uint `json:"number"`
  Upload    uint `json:"upload"`
  Download  uint `json:"download"`
}

type Config struct {
  // the language of the entry point and the path of the entry point: file
  // which will be used as main class running in the simulation; entry point
  // format is [go module]/[nodeStructName]
  Lang string `json:"lang"`
  Entry string `json:"entry"`

  // path to the location of the log.json file
  LogFile string `json:"logFile"`

  // number of peers in the system
  Nodes uint `json:"nodes"`

  // network level topology parameters
  TransitDomains     uint `json:"transitDomains"`
  TransitDomainSize  uint `json:"transitDomainSize"`
  StubDomains        uint `json:"stubDomains"`
  StubDomainSize     uint `json:"stubDomainSize"`

  // interval of virtual time units at which the capacity scheduler runs
  TransferInterval   uint `json:"transferInterval"`

  // list which shows the distribution of node capacities
  CapacityNodes      []NodeConf `json:"capacityNodes"`

  // Latency support - when latency support is off, nodes run directly
  Latency bool `json:"latency"`

  // Run the simulator's event queue with support for parallel events.
  Parallel bool `json:"parallel"`
}
