FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine@sha256:162b257ae1156df92d48e2de5abe69f657fc4961f0ea0e5376837a9ffed8160b AS builder

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

FROM alpine:3.24@sha256:8ddefa941e689fc29abcdeb8dae3b3c6d139cc08ce9a52633931160701770685

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
