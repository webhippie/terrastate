# Terrastate

[![Current Tag](https://img.shields.io/github/v/tag/webhippie/terrastate?sort=semver)](https://github.com/webhippie/terrastate) [![Build Status](https://github.com/webhippie/terrastate/actions/workflows/general.yml/badge.svg)](https://github.com/webhippie/terrastate/actions) [![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org) [![Docker Size](https://img.shields.io/docker/image-size/webhippie/terrastate/latest)](https://hub.docker.com/r/webhippie/terrastate) [![Docker Pulls](https://img.shields.io/docker/pulls/webhippie/terrastate)](https://hub.docker.com/r/webhippie/terrastate) [![Go Reference](https://pkg.go.dev/badge/github.com/webhippie/terrastate.svg)](https://pkg.go.dev/github.com/webhippie/terrastate) [![Go Report Card](https://goreportcard.com/badge/github.com/webhippie/terrastate)](https://goreportcard.com/report/github.com/webhippie/terrastate) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/d2bc4877341f4c7fbf9b4fa62b8d0484)](https://www.codacy.com/gh/webhippie/terrastate/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=webhippie/terrastate&amp;utm_campaign=Badge_Grade)

Terrastate acts as an HTTP backend for Terraform which can store the state
content remotely for you to keep it outside of the repositories containing your
`.tf` files. This is a great alternative if you are not hosting your stuff on
AWS.

## Install

You can download prebuilt binaries from our [GitHub releases][releases], or you
can use our Docker images published on [Docker Hub][dockerhub] or [Quay][quay].
If you need further guidance how to install this take a look at our
[documentation][docs].

## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions][golang]. This project requires
Go >= v1.17, at least that's the version we are using.

```console
git clone https://github.com/webhippie/terrastate.git
cd terrastate

make generate build

./bin/terrastate -h
```

## Security

If you find a security issue please contact
[thomas@webhippie.de](mailto:thomas@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```

[releases]: https://github.com/webhippie/terrastate/releases
[dockerhub]: https://hub.docker.com/r/webhippie/terrastate/tags/
[quay]: https://quay.io/repository/webhippie/terrastate?tab=tags
[docs]: https://webhippie.github.io/terrastate/#getting-started
[golang]: http://golang.org/doc/install.html
