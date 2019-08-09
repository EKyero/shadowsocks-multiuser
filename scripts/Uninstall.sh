#!/usr/bin/env bash
systemctl stop shadowsocks-multiuser
rm -f /etc/systemd/system/shadowsocks-multiuser.service

rm -rf /opt/shadowsocks-multiuser

userdel -r -f shadowsocks-multiuser

exit 0