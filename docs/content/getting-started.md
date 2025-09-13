---
title: "Getting Started"
date: 2022-05-04T00:00:00+00:00
anchor: "getting-started"
weight: 20
---

## Installation

So far we are offering only a few different variants for the installation. You
can choose between [Docker][docker] or pre-built binaries which are stored on
our download mirror and GitHub releases. Maybe we will also provide system
packages for the major distributions later if we see the need for it.

### Docker

Generally we are offering the images through
[quay.io/webhippie/terrastate][quay] and [webhippie/terrastate][dockerhub], so
feel free to choose one of the providers. Maybe we will come up with Kustomize
manifests or some Helm chart.

### Binaries

Simply download a binary matching your operating system and your architecture
from our [downloads][downloads] or the GitHub releases and place it within your
path like `/usr/local/bin` if you are using macOS or Linux.

## Configuration

We provide overall three different variants of configuration. The variant based
on environment variables and commandline flags are split up into global values
and command-specific values.

### Envrionment variables

If you prefer to configure the service with environment variables you can see
the available variables below.

#### Global

TERRASTATE_CONFIG_FILE
: Path to optional config file

TERRASTATE_LOG_LEVEL
: Set logging level, defaults to `info`

TERRASTATE_LOG_COLOR
: Enable colored logging, defaults to `true`

TERRASTATE_LOG_PRETTY
: Enable pretty logging, defaults to `true`

#### Server

TERRASTATE_METRICS_ADDR
: Address to bind the metrics, defaults to `0.0.0.0:8081`

TERRASTATE_METRICS_TOKEN
: Token to make metrics secure

TERRASTATE_SERVER_ADDR
: Address to bind the server, defaults to `0.0.0.0:8080`

TERRASTATE_SERVER_PPROF
: Enable pprof debugging, defaults to `false`

TERRASTATE_SERVER_ROOT
: Root path of the server, defaults to `/`

TERRASTATE_SERVER_HOST
: External access to server, defaults to `http://localhost:8080`

TERRASTATE_SERVER_CERT
: Path to cert for SSL encryption

TERRASTATE_SERVER_KEY
: Path to key for SSL encryption

TERRASTATE_SERVER_STRICT_CURVES
: Use strict SSL curves, defaults to `false`

TERRASTATE_SERVER_STRICT_CIPHERS
: Use strict SSL ciphers, defaults to `false`

TERRASTATE_SERVER_STORAGE
: Folder for storing the states, defaults to `storage/`

TERRASTATE_ENCRYPTION_SECRET
: Secret for file encryption

TERRASTATE_ACCESS_USERNAME
: Username for basic auth

TERRASTATE_ACCESS_PASSWORD
: Password for basic auth

#### Health

TERRASTATE_METRICS_ADDR
: Address to bind the metrics, defaults to `0.0.0.0:8081`

#### State

TERRASTATE_SERVER_STORAGE
: Folder for storing the states, defaults to `storage/`

TERRASTATE_ENCRYPTION_SECRET
: Secret for file encryption

### Commandline flags

If you prefer to configure the service with commandline flags you can see the
available variables below.

#### Global

--config-file
: Path to optional config file

--log-level
: Set logging level, defaults to `info`

--log-color
: Enable colored logging, defaults to `true`

--log-pretty
: Enable pretty logging, defaults to `true`

#### Server

--metrics-addr
: Address to bind the metrics, defaults to `0.0.0.0:8081`

--metrics-token
: Token to make metrics secure

--server-addr
: Address to bind the server, defaults to `0.0.0.0:8080`

--server-pprof
: Enable pprof debugging, defaults to `false`

--server-root
: Root path of the server, defaults to `/`

--server-host
: External access to server, defaults to `http://localhost:8080`

--server-cert
: Path to cert for SSL encryption

--server-key
: Path to key for SSL encryption

--strict-curves
: Use strict SSL curves, defaults to `false`

--strict-ciphers
: Use strict SSL ciphers, defaults to `false`

--storage-path
: Folder for storing the states, defaults to `storage/`

--encryption-secret
: Secret for file encryption

--general-username
: Username for basic auth

--general-password
: Password for basic auth

#### Health

--metrics-addr
: Address to bind the metrics, defaults to `0.0.0.0:8081`

#### State

--storage-path
: Folder for storing the states, defaults to `storage/`

--encryption-secret
: Secret for file encryption

### Configuration file

So far we support multiple file formats like `json`, `yaml`, `hcl` and possibly
even more, if you want to get a full example configuration just take a look at
[our repository][repo], there you can always see the latest configuration
format. These example configs include all available options and the default
values. The configuration file will be automatically loaded if it's placed at
`/etc/terrastate/config.yml`, `${HOME}/.terrastate/config.yml` or
`$(pwd)/terrastate/config.yml`.

## Usage

The program provides a few sub-commands on execution. The available config
methods have already been mentioned above. Generally you can always see a
formated help output if you execute the binary similar to something like
 `terrastate --help`.

Within your Terraform definition you simply got to add this block to get started
using this remote state storage, replace `http://localhost:8080` with your
deployed URL of the Terrastate instance, the rest behind the `/remote` prefix is
entirely up to you to configure as you wish:

{{< highlight yaml >}}
terraform {
  backend "http" {
    address        = "http://localhost:8080/remote/your/state/path"
    lock_address   = "http://localhost:8080/remote/your/state/path"
    unlock_address = "http://localhost:8080/remote/your/state/path"
  }
}
{{< / highlight >}}

[docker]: https://www.docker.com/
[quay]: https://quay.io/repository/webhippie/terrastate
[dockerhub]: https://hub.docker.com/r/webhippie/terrastate
[downloads]: https://dl.webhippie.de/#terrastate/
[repo]: https://github.com/webhippie/terrastate/tree/master/config
