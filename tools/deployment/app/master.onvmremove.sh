#!/bin/bash

# This is setup script for Gamification/Master OnVmRemove.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

# LOADBALANCER="109.231.121.47"
LOADBALANCER_HOST=$(ss-get --timeout 360 loadbalancer.hostname)
LOADBALANCER_API_PORT="1937"

APP_HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
APP_PORT="3000"
APP_TYPE="gamification"

timestamp () {
	date +"%F %T,%3N"
}

inform_loadbalancer () {
	retries=0
	until curl -X DELETE http://$LOADBALANCER_HOST:$LOADBALANCER_API_PORT/$APP_TYPE/$APP_HOST:$APP_PORT; do
		echo "[$(timestamp)] Load Balancer is not running, retry... [$(( retries++ ))]"
		sleep 5
	done
}

echo "[$(timestamp)] ---- 1. Inform Load Balancer (Remove) ----"
inform_loadbalancer

echo "[$(timestamp)] ---- Completed ----"

exit 0
