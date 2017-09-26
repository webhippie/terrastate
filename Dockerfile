FROM alpine:edge
MAINTAINER Thomas Boerger <thomas@webhippie.de>

EXPOSE 8080
VOLUME ["/var/lib/terrastate"]

RUN apk update && \
  apk add \
    ca-certificates \
    bash && \
  rm -rf \
    /var/cache/apk/* && \
  addgroup \
    -g 1000 \
    terrastate && \
  adduser -D \
    -h /var/lib/terrastate \
    -s /bin/bash \
    -G terrastate \
    -u 1000 \
    terrastate

COPY terrastate /usr/bin/

ENV TERRASTATE_SERVER_STORAGE /var/lib/terrastate

USER terrastate
ENTRYPOINT ["/usr/bin/terrastate"]
CMD ["server"]
