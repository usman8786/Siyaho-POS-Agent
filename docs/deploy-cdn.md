# CDN deploy — pos.siyaho.com

Agent installers are served from:

```
https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe
https://pos.siyaho.com/downloads/print-agent/releases.json
```

Files live **outside** `/var/deployments/pos-web` so POS web `rsync --delete` does not remove them.

## Server path

```
/var/deployments/pos-downloads/print-agent/
```

## Server setup (one-time)

```bash
sudo mkdir -p /var/deployments/pos-downloads/print-agent
sudo chown -R siyaho:siyaho /var/deployments/pos-downloads
```

Update nginx from `Siyaho-POS-Web/deploy/nginx-pos-web.conf`, then:

```bash
sudo nginx -t && sudo systemctl reload nginx
```

## Manual upload (first time or recovery)

From your PC (after downloading from GitHub Releases):

```bash
scp SiyahoPrinterAgent-Setup-1.0.0.exe siyaho@103.18.23.182:/var/deployments/pos-downloads/print-agent/
scp releases.json siyaho@103.18.23.182:/var/deployments/pos-downloads/print-agent/
ssh siyaho@103.18.23.182 "cd /var/deployments/pos-downloads/print-agent && ln -sf SiyahoPrinterAgent-Setup-1.0.0.exe SiyahoPrinterAgent-Setup-latest.exe"
```

## CI

The GitHub release workflow uploads to this path on tag `v*`.
