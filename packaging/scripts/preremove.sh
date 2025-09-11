#!/bin/sh
set -e

systemctl stop terrastate.service || true
systemctl disable terrastate.service || true
