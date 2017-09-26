# Terrastate

[![Build Status](http://github.dronehippie.de/api/badges/webhippie/terrastate/status.svg)](http://github.dronehippie.de/webhippie/terrastate)
[![Go Doc](https://godoc.org/github.com/webhippie/terrastate?status.svg)](http://godoc.org/github.com/webhippie/terrastate)
[![Go Report](https://goreportcard.com/badge/github.com/webhippie/terrastate)](https://goreportcard.com/report/github.com/webhippie/terrastate)
[![Sourcegraph](https://sourcegraph.com/github.com/webhippie/terrastate/-/badge.svg)](https://sourcegraph.com/github.com/webhippie/terrastate?badge)
[![](https://images.microbadger.com/badges/image/tboerger/terrastate.svg)](http://microbadger.com/images/tboerger/terrastate "Get your own image badge on microbadger.com")
[![Join the chat at https://gitter.im/webhippie/general](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/webhippie/general)
[![Stories in Ready](https://badge.waffle.io/webhippie/terrastate.svg?label=ready&title=Ready)](http://waffle.io/webhippie/terrastate)

**This project is under heavy development, it's not in a working state yet!**

Terrastate acts as an HTTP backend for Terraform which can store the state
content remotely for you to keep it outside of the repositories containing your
`.tf` files. This is a great alternative if you are not hosting your stuff on
AWS.


## Install

You can download prebuilt binaries from the GitHub releases or from our
[download site](http://dl.webhippie.de/misc/terrastate). You are a Mac user?
Just take a look at our [homebrew formula](https://github.com/webhippie/homebrew-webhippie).
If you are missing an architecture just write us on our nice
[Gitter](https://gitter.im/webhippie/general) chat. If you find a security issue
please contact thomas@webhippie.de first.


## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions](http://golang.org/doc/install.html).
As this project relies on vendoring of the dependencies and we are not
exporting `GO15VENDOREXPERIMENT=1` within our makefile you have to use a Go
version `>= 1.6`. It is also possible to just simply execute the
`go get github.com/webhippie/terrastate/cmd/terrastate` command, but we
prefer to use our `Makefile`:

```bash
go get -d github.com/webhippie/terrastate
cd $GOPATH/src/github.com/webhippie/terrastate
make clean build

./terrastate -h
```


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2017 Thomas Boerger <http://www.webhippie.de>
```
