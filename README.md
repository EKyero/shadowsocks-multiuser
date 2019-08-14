# shadowsocks-multiuser
**注意，此后端仅支持原版 SS 而且加密方式较少！**

## 支持的加密方式
```
AEAD_AES_128_GCM (AES-128-GCM)
AEAD_AES_192_GCM (AES-192-GCM)
AEAD_AES_256_GCM (AES-256-GCM)
AEAD_CHACHA20_POLY1305 (CHACHA20-POLY1305)
AEAD_XCHACHA20_POLY1305 (XCHACHA20-POLY1305)
RC4-MD5
AES-128-CFB
AES-192-CFB
AES-256-CFB
AES-128-CTR
AES-192-CTR
AES-256-CTR
CHACHA20
CHACHA20-IETF
XCHACHA20
```

## 命令行参数
```
Usage of shadowsocks-multiuser:
  -dbhost string
        Database hostname (default "localhost")
  -dbname string
        Database name (default "sspanel")
  -dbpass string
        Database password (default "123456")
  -dbport int
        Database port (default 3306)
  -dbuser string
        Database username (default "root")
  -listcipher
        List all cipher
  -nodeid int
        Node ID (default -1)
  -syncinterval int
        Sync interval (default 30)
  -udp
        UDP forward
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