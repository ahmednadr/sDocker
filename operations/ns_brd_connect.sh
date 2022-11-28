NS="ns1"
VETH="veth1"
VPEER="vpeer1"
VPEER_ADDR="10.10.0.10"

# remove namespace if it exists.
ip netns del $NS &>/dev/null

# create veth link
ip link add ${VETH} type veth peer name ${VPEER}

# setup veth link
ip link set ${VETH} up

# add peers to ns
ip link set ${VPEER} netns ${NS}

# setup loopback interface
ip netns exec ${NS} ip link set lo up

# assign ip address to ns interfaces
ip netns exec ${NS} ip addr add ${VPEER_ADDR}/16 dev ${VPEER}

# assign veth pairs to bridge
ip link set ${VETH} master ${BR_DEV}

# add default routes for ns
ip netns exec ${NS} ip route add default via ${BR_ADDR}

