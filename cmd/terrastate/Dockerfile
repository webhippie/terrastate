FROM --platform=$BUILDPLATFORM golang:1.25.4-alpine3.21@sha256:1b84283ebeef726bc5fa9fec8deb36828aabfa12fe3a28b4fb0a4b2eafafe38c AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/terrastate
COPY . /go/src/terrastate/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 terrastate && \
    adduser -D -u 1337 -h /var/lib/terrastate -G terrastate terrastate

EXPOSE 8080 8081
VOLUME ["/var/lib/terrastate"]
ENTRYPOINT ["/usr/bin/terrastate"]
CMD ["server"]
HEALTHCHECK CMD ["/usr/bin/terrastate", "health"]

ENV TERRASTATE_STORAGE=/var/lib/terrastate

COPY --from=builder /go/src/terrastate/bin/terrastate /usr/bin/terrastate
WORKDIR /var/lib/terrastate
USER terrastate
