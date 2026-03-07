# Security Analysis - ssh.joledev.com

## Architecture

The SSH server is built with [Charm Wish](https://github.com/charmbracelet/wish), which implements
a custom SSH server in Go using `gliderlabs/ssh` (not OpenSSH). It only serves a TUI application
and never exposes a shell.

## Threat Model

### What this app does:
- Listens on a TCP port and accepts SSH connections
- Authenticates all connections (no credentials required - it's a public portfolio)
- Renders a Bubbletea TUI and streams it to the client
- Reads local files (posts, songs) in read-only mode

### What this app does NOT do:
- Does NOT expose a shell or allow command execution
- Does NOT accept file uploads
- Does NOT store user data or credentials
- Does NOT run as root
- Does NOT modify any files on disk

## Risk Assessment

### LOW RISK
| Risk | Mitigation |
|------|-----------|
| Port exposure | Only the portfolio port (2222) is exposed, not the admin SSH |
| Resource exhaustion | Wish handles connection limits; systemd can enforce memory/CPU limits |
| SSH key theft | Server keys stored in restricted directory (700 permissions) |
| Code injection | No user input is executed; TUI only accepts navigation keypresses |

### MEDIUM RISK
| Risk | Mitigation |
|------|-----------|
| DoS via many connections | Add `LimitNOFILE=1024` and `MemoryMax=256M` to systemd service |
| Wish library vulnerability | Keep dependencies updated; `go get -u` regularly |
| Information disclosure | Repo is public - NO secrets in code, use env vars for config |

### RECOMMENDATIONS
1. **Run on a dedicated port (2222)**, not port 22 - keep your admin SSH on a separate port
2. **Dedicated user** with no shell (`/usr/sbin/nologin`) and minimal permissions
3. **systemd hardening** - the service file includes `NoNewPrivileges`, `ProtectSystem=strict`, etc.
4. **Never commit** `.ssh/` directory, `.env` files, or any credentials
5. **Keep Go dependencies updated** monthly
6. **Use fail2ban** on your admin SSH port (not needed for the portfolio port)
7. **Firewall** - only open the portfolio port and your admin SSH port

## Port Strategy

```
Port 22    -> iptables REDIRECT -> port 2223 (portfolio app)
Port 2222  -> admin SSH (used by CI/CD)
Port 2223  -> portfolio app (Wish SSH server, direct access)
```

Users simply run: `ssh ssh.joledev.com` (port 22, redirected to portfolio)
Admin access: `ssh -p 2222 user@your-server`

The iptables NAT rule:
```
iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-port 2223
```

## Public Repo Checklist
- [ ] `.gitignore` includes `.ssh/`, `.env`, `*.pem`, `*.key`
- [ ] No hardcoded IPs, passwords, or tokens
- [ ] Server config via environment variables only
- [ ] No sensitive file paths exposed
