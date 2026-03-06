FROM --platform=$BUILDPLATFORM golang:1.26.1-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS builder

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

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

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
