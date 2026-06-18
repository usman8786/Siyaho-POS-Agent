# CDN deploy — pos.siyaho.com

Agent installers are served from:

```
https://pos.siyaho.com/downloads/print-agent/SiyahoPrinterAgent-Setup-latest.exe
https://pos.siyaho.com/downloads/print-agent/releases.json
```

## Server setup (one-time)

SSH to `103.18.23.182`:

```bash
sudo mkdir -p /var/deployments/pos-web/downloads/print-agent
sudo chown -R $USER:$USER /var/deployments/pos-web/downloads
```

Update nginx (`Siyaho-POS-Web/deploy/nginx-pos-web.conf`) with the `/downloads/print-agent/` location, then:

```bash
sudo nginx -t && sudo systemctl reload nginx
```

## Manual upload

```bash
scp dist/SiyahoPrinterAgent-Setup-1.0.0.exe root@103.18.23.182:/var/deployments/pos-web/downloads/print-agent/
scp releases.json root@103.18.23.182:/var/deployments/pos-web/downloads/print-agent/
ssh root@103.18.23.182 "cd /var/deployments/pos-web/downloads/print-agent && ln -sf SiyahoPrinterAgent-Setup-1.0.0.exe SiyahoPrinterAgent-Setup-latest.exe"
```

## CI

The GitHub release workflow uploads to CDN when `SELFHOSTED_SSH_PRIVATE_KEY` is configured (same secret as POS web deploy).
