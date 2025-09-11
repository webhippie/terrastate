#!/bin/sh
set -e

if ! getent group terrastate >/dev/null 2>&1; then
    groupadd --system terrastate
fi

if ! getent passwd terrastate >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/terrastate --shell /bin/bash -g terrastate terrastate
fi
