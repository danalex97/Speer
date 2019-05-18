<img src="docs/pics/logo.png" width="300">

[![Build Status](https://travis-ci.org/danalex97/Speer.svg?branch=master)](https://travis-ci.org/danalex97/Speer) [![Coverage Status](https://coveralls.io/repos/github/danalex97/Speer/badge.svg?branch=master)](https://coveralls.io/github/danalex97/Speer?branch=master)
[![GoDoc](https://godoc.org/github.com/danalex97/Speer?status.png)](https://godoc.org/github.com/danalex97/Speer)


A discrete event **S**imulator for **peer**-to-peer network modeling. **Speer is made for students, researchers and hobbyists.** It's goal is to allow
them to easily implement, simulate and study peer to peer networks.

It combines event-driven simulations with cycle-based concepts and allows parallelization by taking advantage of Go’s concurrency features.

## Quickstart

After getting **Golang >= 1.6** and setting **$GOPATH**, install Speer via:
```
curl https://raw.githubusercontent.com/danalex97/Speer/master/install.sh | bash
```

Now, you can run a simulation from a JSON configuration as follows:
```
speer -config=[path_to_configuration]
```

The see the options provided by `speer.go` run:
```
speer -h
```

You can run the default example in `examples/broadcast.go` via the command:
```
speer
```

## Table of contents

- [Motivation & FAQ](docs/motivation.md)
- [Architecture](docs/architecture.md)
  - [Event simulator](docs/events.md)
  - [Latency](docs/latency.md)
  - [Capacity](docs/capacity.md)
  - [Optimizations](docs/optimizations.md)
- User guide
  - [Interfaces](docs/interfaces.md)
  - [Examples](docs/examples.md)
  - [Running a simulation](docs/running.md)

## How to contribute

- [Contribution guide](.github/CONTRIBUTING.md)
- [Roadmap](docs/roadmap.md)

## Projects using Speer

- [CacheTorrent](https://github.com/danalex97/nfsTorrent) - is a file sharing system based on leader election, caches and indirect requests

<img src="https://raw.githubusercontent.com/danalex97/nfsTorrent/master/docs/pics/cache.png" width="500">
