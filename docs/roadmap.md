## Roadmap

The direction in which the project goes should prioritize the following concerns:
 - developer support
 - framework usability
 - simulation accuracy

*TL;DR: Check version 1.2 for the current prioritized work.*

The Roadmap below is not fixed and contributions for parts in letter stages can be done earlier. If as a contributor you are interested in working on a part planned later, you are free to do so and your code will be reviewed!

#### Version 1.0 [done]

- no configuration files for running simulations
- no good examples
- no main runner with specified entry point entry point

#### Version 1.1 [done]

- clean code-base
- support for configuration files
- reviewing standards specified
- at least 2 good usage examples
- main command line utility which allows definition of:
  - path for user-supplied code (i.e. a *.go* file)
  - path for configuration file
- setup of a front-end for visualization in a JS framework(e.g. [React](https://reactjs.org/) & [D3.js](https://d3js.org/))

#### Version 1.2 [intent: September 2019]

- basic front-end with visualization tools for:
  - hop count distribution
  - latency distribution
  - throughput distribution
  - network layer topology
  - transport layer topology
- progress and safety support
- handshake support
- failure support
  - node failures
  - joins and leaves
  - lossy links
- routing algorithms
  - using [BGP](https://en.wikipedia.org/wiki/Border_Gateway_Protocol) together with OSPF
  - support for dynamic latencies

#### Version 1.3
- non-preemptive Python SDK implementation:
    - possibility user supplied code in Python
    - Golang examples in Python as well
    - library wrapper and deployment for pip
- non-preemptive Rust SDK implementation

#### Version 2.0

During start of version 2.0 we want to support requirements by universities interested in using Speer in distributed algorithms courses. Things we consider are related to usability and could include:
- implementing new SDKs
- changing the front-end
- changing the setup process
- changing programming primitives

#### Version 2.1

This work is mainly focused on improving the simulations themselves. Some examples include:
- moving flow simulations to network level(see this [paper](https://dl.acm.org/citation.cfm?id=1272986))
- realistic topologies: by using real datasets(e.g. [Archipelago Measurement Infrastructure](http://www.caida.org/projects/ark/))

Other work can be related to improving the current simulation infrastructure. Some examples are:
- replacing the heap with O(1) lazy structures
- parallelize the flow allocation algorithms
- faster ICP for SDKs

#### Version 3.0

- infrastructure to allow deployment of multiple simulations on multiple machines(each on 1 machine)
- infrastructure to allow deployment of a single simulation on multiple machines
