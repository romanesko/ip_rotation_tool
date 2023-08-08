### run

expecting `conf.yaml` in the same directory

```yaml
port: 8080
endpoint: "/some/endpoint"
ips:
  - 99.84.140.147
  - 143.204.79.125
  - 54.240.188.143
  - 143.204.127.42
  - 13.35.51.41
  - 13.35.55.41
  - 99.84.58.138
```

### quick deploy

```bash
#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -o ip_rotation_tool
scp ip_rotation_tool conf.yaml <server>:ip_rotation_tool
```

### supervisord

1. install supervisord
```bash
sudo apt-get update
sudo apt-get install -y supervisor 
sudo service supervisor start
```

2. add config

```bash
sudo vim /etc/supervisor/conf.d/ip_rotation_tool.conf
```

with content

```bash
[program:ip_rotation_tool]
directory=/root/home
command=/root/home/ip_rotation_tool
autostart=true
autorestart=true
stderr_logfile=/var/log/api.err
stdout_logfile=/var/log/api.log
```

3. restart supervisor

```bash
sudo supervisorctl reload
```

4. check the process status:

```bash
sudo supervisorctl status
```