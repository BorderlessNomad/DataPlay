#!/bin/bash

# This is setup script for PostgreSQL OnVmRemove.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

#PGPOOL_API_HOST="109.231.124.33"
PGPOOL_API_HOST=$(ss-get --timeout 360 pgpool.hostname)
PGPOOL_API_PORT="1937"

APP_HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')

timestamp () {
	date +"%F %T,%3N"
}

inform_pgpool () {
	retries=0
	until curl -X DELETE http://$PGPOOL_API_HOST:$PGPOOL_API_PORT/$APP_HOST; do
		echo "[$(timestamp)] PGPOOL Server is not up yet, retry... [$(( retries++ ))]"
		sleep 5
	done
}

echo "[$(timestamp)] ---- 1. Inform pgpool (Remove) ----"
inform_pgpool

echo "[$(timestamp)] ---- Completed ----"

exit 0
