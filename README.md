# shadowsocks-multiuser
**注意，此后端仅支持原版 SS 而且加密方式较少！**

## 支持的加密方式
```
AES-128-GCM
AES-192-GCM
AES-256-GCM
XCHACHA20
XCHACHA20-POLY1305
AES-128-CFB
AES-192-CFB
AES-256-CFB
AES-128-CTR
AES-192-CTR
AES-256-CTR
CHACHA20
CHACHA20-IETF
CHACHA20-POLY1305
RC4-MD5
```

## Ubuntu 自动安装脚本
1. 运行脚本
```bash
curl -fsSL https://raw.githubusercontent.com/NetchX/shadowsocks-multiuser/master/scripts/Ubuntu.sh | bash
```

2. 编辑配置文件调整参数
```bash
vim /etc/systemd/system/shadowsocks-multiuser.service
```

3. 重启 Systemd 并开启服务设置自启
```bash
systemctl daemon-reload
systemctl start shadowsocks-multiuser
systemctl enable shadowsocks-multiuser
systemctl status shadowsocks-multiuser
```