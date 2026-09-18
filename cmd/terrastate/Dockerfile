FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine@sha256:e9bbdf282b51ac8b34c46e5f31d2d56e7bad60366c35f08d2f295b921b13388b AS builder

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

FROM alpine:3.24@sha256:e7c4abb69531cb09e2a2bbb56fad3367ab694865c49df898c1c683185cc4376c

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
