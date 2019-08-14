#!/usr/bin/env bash

# Stop service
systemctl stop shadowsocks-multiuser

# Remove service
rm -f /etc/systemd/system/shadowsocks-multiuser.service

# Remove files
rm -rf /opt/shadowsocks-multiuser

# Return 0
exit 0