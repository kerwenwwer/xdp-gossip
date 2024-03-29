#!/bin/bash

# Usage: veth.sh <container name>
# Description: Get the veth interface of a container

NAME=$1
PID=$(docker inspect $NAME --format "{{.State.Pid}}")
while read iface id; do
    [[ "$iface" == lo ]] && continue
    veth=$(ip -br addr | sed -nre "s/(veth.*)@if$id.*/\1/p")
    echo -e "$NAME $iface $veth"
done < <(</proc/$PID/net/igmp awk '/^[0-9]+/{print $2 " " $1;}')
