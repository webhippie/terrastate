#!/bin/sh
set -e

chown -R terrastate:terrastate /etc/terrastate
chown -R terrastate:terrastate /var/lib/terrastate
chmod 750 /var/lib/terrastate

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet terrastate.service; then
        systemctl restart terrastate.service
    fi
fi
