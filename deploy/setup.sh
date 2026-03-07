#!/bin/bash
# Setup script for ssh.joledev.com on a VPS
# Admin SSH is on port 2222, portfolio app on port 2223
# iptables redirects port 22 -> 2223 for the portfolio
#
# Run as root or with sudo

set -euo pipefail

APP_USER="sshjoledev"
APP_DIR="/opt/ssh-joledev"
APP_PORT=2223

echo "==> Creating dedicated user (no shell, no home)..."
useradd --system --no-create-home --shell /usr/sbin/nologin "$APP_USER" 2>/dev/null || true

echo "==> Creating app directory..."
mkdir -p "$APP_DIR"/{.ssh,data,posts/es,posts/en}

echo "==> Copying files..."
cp ssh.joledev "$APP_DIR/"
cp -r data/* "$APP_DIR/data/" 2>/dev/null || true
cp -r posts/* "$APP_DIR/posts/" 2>/dev/null || true

echo "==> Setting permissions..."
chown -R "$APP_USER:$APP_USER" "$APP_DIR"
chmod 700 "$APP_DIR/.ssh"
chmod 755 "$APP_DIR/ssh.joledev"

echo "==> Installing systemd service..."
cp deploy/ssh-joledev.service /etc/systemd/system/ssh-joledev.service
systemctl daemon-reload
systemctl enable ssh-joledev
systemctl restart ssh-joledev

echo "==> Setting up iptables: port 22 -> $APP_PORT..."
# Remove old rule if exists
iptables -t nat -D PREROUTING -p tcp --dport 22 -j REDIRECT --to-port "$APP_PORT" 2>/dev/null || true
# Add redirect rule
iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-port "$APP_PORT"

# Make iptables rule persistent
if command -v netfilter-persistent &>/dev/null; then
    netfilter-persistent save
elif [ -f /etc/iptables/rules.v4 ]; then
    iptables-save > /etc/iptables/rules.v4
else
    echo "    WARNING: Install iptables-persistent to make the rule survive reboots:"
    echo "    apt install iptables-persistent"
fi

# Open port 22 in firewall if needed
if command -v ufw &>/dev/null; then
    ufw allow 22/tcp 2>/dev/null || true
    ufw allow "$APP_PORT/tcp" 2>/dev/null || true
fi

echo ""
echo "==> Done!"
echo "    Portfolio app: port $APP_PORT (direct)"
echo "    Port 22 redirects to $APP_PORT via iptables"
echo "    Admin SSH: port 2222 (unchanged)"
echo ""
echo "    Test: ssh ssh.joledev.com"
echo "    Admin: ssh -p 2222 user@ssh.joledev.com"
