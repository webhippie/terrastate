FROM webhippie/alpine:latest AS build
RUN apk add --no-cache ca-certificates mailcap

FROM scratch

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="Terrastate" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

EXPOSE 8080
VOLUME ["/var/lib/terrastate"]

ENV TERRASTATE_STORAGE /var/lib/terrastate

ENTRYPOINT ["/usr/bin/terrastate"]
CMD ["server"]

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/mime.types /etc/

COPY dist/binaries/terrastate-*-linux-arm-5 /usr/bin/terrastate
