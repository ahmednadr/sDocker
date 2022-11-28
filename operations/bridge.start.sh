#!/usr/bin/env bash

if [[ $EUID -ne 0 ]]; then
    echo "You must be root to run this script"
    exit 1
fi


BR_ADDR="10.10.0.1"
BR_DEV="br0"

# setup bridge
ip link add ${BR_DEV} type bridge
ip link set ${BR_DEV} up