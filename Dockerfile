FROM alpine:edge
MAINTAINER Thomas Boerger <thomas@webhippie.de>

LABEL org.label-schema.vcs-url="https://github.com/webhippie/terrastate.git"
LABEL org.label-schema.name="Terrastate"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"

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

ENV TERRASTATE_SERVER_STORAGE /var/lib/terrastate

USER terrastate
ENTRYPOINT ["/usr/bin/terrastate"]
CMD ["server"]

COPY terrastate /usr/bin/
