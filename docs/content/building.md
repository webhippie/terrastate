---
title: "Building"
date: 2022-05-03T00:00:00+00:00
anchor: "building"
weight: 30
---

As this project is built with Go you need to install Go first. The installation
of Go is out of the scope of this document, please follow the
[official documentation][golang]. After the installation of Go you need to get
the sources:

{{< highlight txt >}}
git clone https://github.com/webhippie/terrastate.git
cd terrastate/
{{< / highlight >}}

All required tool besides Go itself are bundled by Go modules, all you need is
part of the `Makfile`:

{{< highlight txt >}}
make generate build
{{< / highlight >}}

Finally you should have the binary within the `bin/` folder now, give it a try
with `./bin/terrastate -h` to see all available options.

[golang]: https://golang.org/doc/install
