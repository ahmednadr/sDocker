#!/usr/bin/env bash

if [[ $EUID -ne 0 ]]; then
    echo "You must be root to run this script"
    exit 1
fi

# check if bridge already exists

BR_ADDR="10.10.0.1"
BR_DEV="sDocker0"

# setup bridge
ip link add ${BR_DEV} type bridge
ip link set ${BR_DEV} up

# setup bridge ip
ip addr add ${BR_ADDR}/16 dev ${BR_DEV}

# enable ip forwarding
bash -c 'echo 1 > /proc/sys/net/ipv4/ip_forward'

# Flush nat rules.
iptables -t nat -F


iptables -t nat -A POSTROUTING -s ${BR_ADDR}/16 ! -o ${BR_DEV} -j MASQUERADE