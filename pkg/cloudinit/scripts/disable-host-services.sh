#!/bin/bash -xe

# Usage:
#   $0
#
# Assumptions:
#   - systemctl is available

for svc in kubelet containerd; do
  systemctl stop "${svc}" || true
  systemctl disable "${svc}" || true
  ln -sf /dev/null "/etc/systemd/system/${svc}.service"
done
systemctl daemon-reload
