#!/bin/bash

# This is setup script for Gamification/Master OnVmRemove.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

# LOADBALANCER="109.231.121.47"
LOADBALANCER_HOST=$(ss-get --timeout 360 loadbalancer.hostname)
LOADBALANCER_PORT="1937"
HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
PORT="3000"

timestamp () {
	date +"%F %T,%3N"
}

inform_loadbalancer () {
	PAYLOAD="$HOST:$PORT"
	curl -i -X DELETE http://$LOADBALANCER_HOST:$LOADBALANCER_PORT/$PAYLOAD
}

echo "[$(timestamp)] ---- 1. Inform Load Balancer (Remove) ----"
inform_loadbalancer

echo "[$(timestamp)] ---- Completed ----"

exit 0
