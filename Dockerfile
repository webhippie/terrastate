FROM webhippie/alpine:latest

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="Terrastate" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

EXPOSE 8080
VOLUME ["/var/lib/terrastate"]

ENV TERRASTATE_STORAGE /var/lib/terrastate

ENTRYPOINT ["/usr/bin/terrastate"]
CMD ["server"]

RUN apk add --no-cache ca-certificates mailcap bash

COPY dist/binaries/terrastate-*-linux-amd64 /usr/bin/terrastate
