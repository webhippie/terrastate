# Terrastate

[![Build Status](http://github.dronehippie.de/api/badges/webhippie/terrastate/status.svg)](http://github.dronehippie.de/webhippie/terrastate)
[![Stories in Ready](https://badge.waffle.io/webhippie/terrastate.svg?label=ready&title=Ready)](http://waffle.io/webhippie/terrastate)
[![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie%3Amatrix.org-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d2bc4877341f4c7fbf9b4fa62b8d0484)](https://www.codacy.com/app/webhippie/terrastate?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=webhippie/terrastate&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/webhippie/terrastate?status.svg)](http://godoc.org/github.com/webhippie/terrastate)
[![Go Report](https://goreportcard.com/badge/github.com/webhippie/terrastate)](https://goreportcard.com/report/github.com/webhippie/terrastate)
[![](https://images.microbadger.com/badges/image/tboerger/terrastate.svg)](http://microbadger.com/images/tboerger/terrastate "Get your own image badge on microbadger.com")
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/1828/badge)](https://bestpractices.coreinfrastructure.org/projects/1828)

Terrastate acts as an HTTP backend for Terraform which can store the state content remotely for you to keep it outside of the repositories containing your `.tf` files. This is a great alternative if you are not hosting your stuff on AWS.


## Docs

Our documentation gets generated directly out of the [docs/](docs/) folder, it get's built via Drone and published to GitHub pages. You can find the documentation at [https://webhippie.github.io/terrastate/](https://webhippie.github.io/terrastate/).


## Install

You can download prebuilt binaries from the GitHub releases or from our [download site](https://dl.webhippie.de/terrastate/master/). You are a Mac user? Just take a look at our [homebrew formula](https://github.com/webhippie/homebrew-webhippie).


## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.8.

```bash
go get -d github.com/webhippie/terrastate
cd $GOPATH/src/github.com/webhippie/terrastate
make clean generate build

./bin/terrastate -h
```


## Security

If you find a security issue please contact thomas@webhippie.de first.


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2018 Thomas Boerger <http://www.webhippie.de>
```
