#!/usr/bin/env bash
# Make sure the script runs with super user privileges.
[ "$UID" -eq 0 ] || exec sudo bash "$0" "$@"

modprobe vcan
ip link add dev vcan0 type vcan
ip link set up vcan0
ip link show vcan0
