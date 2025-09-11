#!/bin/sh
set -e

if [ ! -d /var/lib/terrastate ] && [ ! -d /etc/terrastate ]; then
    userdel terrastate 2>/dev/null || true
    groupdel terrastate 2>/dev/null || true
fi
