# Poor Man's exe.dev - SSH API Reference

The primary way to interact with the platform is via SSH.

## Basic Commands

### List VMs
```bash
ssh poor-exe.yourdomain.com ls
```
To get script-friendly JSON:
```bash
ssh poor-exe.yourdomain.com ls --json | jq .
```

### Create a VM
```bash
ssh poor-exe.yourdomain.com new --name=bloggy --image=nginx:alpine
```
Returns endpoints and connection details.

### Delete a VM
```bash
ssh poor-exe.yourdomain.com rm bloggy
```

### User Info
```bash
ssh poor-exe.yourdomain.com whoami
```

---

## Share Commands

Manage public access and port mapping.

### Make Public (No login required for HTTP)
```bash
ssh poor-exe.yourdomain.com share set-public bloggy
```

### Make Private (Default)
```bash
ssh poor-exe.yourdomain.com share set-private bloggy
```

### Map HTTP Port
Change which internal port is exposed to the public internet (default: 80).
```bash
ssh poor-exe.yourdomain.com share port bloggy 8080
```

### Management via Email
*Note: In this MVP, this adds users to the allowlist for private VMs.*
```bash
ssh poor-exe.yourdomain.com share add bloggy friend@example.com
ssh poor-exe.yourdomain.com share remove bloggy friend@example.com
```

---

## Connecting to VMs

Once created, you can SSH directly into the VM shell:

```bash
ssh bloggy@poor-exe.yourdomain.com
```

### ssh_config optimization
Add this to your `~/.ssh/config` for easier access:

```ssh-config
Host *.yourdomain.com
  HostName poor-exe.yourdomain.com
  User %h
```

Now you can just run:
```bash
ssh bloggy.yourdomain.com
```
