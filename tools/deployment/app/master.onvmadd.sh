#!/bin/bash

# This is setup script for Gamification/Master OnVmAdd.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

# LOADBALANCER="109.231.121.47"
LOADBALANCER_HOST=$(ss-get --timeout 360 loadbalancer.hostname)
LOADBALANCER_PORT="1937"
HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
PORT="80"

timestamp () {
	date +"%F %T,%3N"
}

inform_loadbalancer () {
	PAYLOAD={\"ip\":\"$HOST:$PORT\"}
	curl -i -H "Content-Type: application/json" -X POST -d '$PAYLOAD' http://$LOADBALANCER_HOST:$LOADBALANCER_PORT
}

echo "[$(timestamp)] ---- 1. Inform Load Balancer (Add) ----"
inform_loadbalancer

echo "[$(timestamp)] ---- Completed ----"

exit 0
