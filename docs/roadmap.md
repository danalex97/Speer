## Roadmap

The direction in which the project goes should prioritize the following concerns:
 - developer support
 - framework usability
 - simulation accuracy

*TL;DR: Check version 1.1 for the current prioritized work.*

We decide to prioritize **developer support** in order to be able to acquire a bigger developer community working on Speer. Moreover, we want to acquire **framework usability** before **simulation accuracy** so that our software is usable for a winder community, including university courses. We encountered some interest from universities for using the tool in the distributed algorithms courses, hence the prioritization of framework usability.

#### Workload and contribution

We aim to split the work for the project to be **small over larger periods of times**. We aim for a mean contribution time between 2-3 hours weekly.

We aim for the contributor work onto this project to be recreative rather than stressful. See our [contributor guidelines](contribute.md) for more details.

The Roadmap below is not fixed and contributions for parts in letter stages can be done earlier. If as a contributor you are interested in working on a part planned later, you are free to do so and your code will be reviewed!

#### Version 1.0 [current]

Currently, the system has a few notable problems in terms of developer support and usability:
- no configuration files for running simulations
- no good examples
- no main runner with specified entry point entry point

#### Version 1.1 [intent: December 2018]

We aim for this version to be **contributor friendly**:
- clean code-base
- support for **configuration files**
- reviewing standards specified and established core reviewers
- at least 2 good usage examples
- main command line utility which allows definition of:
  - **path for user-supplied code** (i.e. a *.go* file)
  - path for configuration file
- setup of a **front-end for visualization** in a JS framework(e.g. [React](https://reactjs.org/) & [D3.js](https://d3js.org/))

#### Version 1.2 [intent: March/April 2019]

At this version we want to **expand our contributor base**. The main purpose of this stage is to work on **usability**:
- **basic front-end** with visualization tools for:
  - hop count distribution
  - latency distribution
  - throughput distribution
  - network layer topology
  - transport layer topology
- non-preemptive **Python SDK** implementation:
  - possibility user supplied code in Python
  - Golang examples in Python as well
  - library wrapper and deployment for pip

#### Version 2.0 [intent: September 2019]

During start of version 2.0 we want to already be able to start **talks with universities for using Speer in distributed algorithms courses.**

The main focus of this stage is **gathering requirements** from universities and implementing them into Speer. Things we consider are related to usability and could include:
- implementing new SDKs
- changing the front-end
- changing the setup process
- changing programming primitives

The work planned during this period can intersect with parts planned for the next version.

#### Version 2.1 [intent: End of 2019/Beginning of 2020]

This work is mainly focused on improving the simulations themselves, in particular:
- moving flow simulations to network level(see this [paper](https://dl.acm.org/citation.cfm?id=1272986))
- dynamic topologies(initially model failures as Possion processes)
- dynamic workloads(model node departure and arrival)
- node failure(model node failure - hard and soft)
- using [BGP](https://en.wikipedia.org/wiki/Border_Gateway_Protocol) together with OSPF
- using [B4](https://dl.acm.org/citation.cfm?id=2486019) for WANs
- realistic topologies: by using real datasets(e.g. [Archipelago Measurement Infrastructure](http://www.caida.org/projects/ark/))

Other work can be related to improving the current simulation infrastructure. Some examples are:
- replacing the heap with O(1) lazy structures
- parallelize the flow allocation algorithms
- use memory mapped files to provide faster ICP for Python SDK

We are open to new ideas besides the ones described above.

#### Version 3.0 [intent: sometime in the distant future]

The final goal of the simulator is to allow for **large scale simulations**. This can be done by allowing deployment on multiple machines.

The two main lines of work are:
- infrastructure to allow deployment of multiple simulations on multiple machines(each on 1 machine)
- infrastructure to allow deployment of a single simulation on multiple machines
