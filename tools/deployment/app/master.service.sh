#!/bin/bash

# This is service script for Master instance.
# Used for updating app on server after deployment.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

DEST="/home/ubuntu/www"
APP="dataplay"

APP_HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
APP_PORT="3000"
APP_TYPE="master"
APP_MODE="3"

. /etc/profile

# LOADBALANCER_HOST="109.231.121.26"
LOADBALANCER_HOST="$DP_LOADBALANCER_HOST"
LOADBALANCER_REQUEST_PORT="$DP_LOADBALANCER_REQUEST_PORT"
LOADBALANCER_API_PORT="$DP_LOADBALANCER_API_PORT"

timestamp () {
	date +"%F %T,%3N"
}

run_master_server () {
	URL="https://github.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO"

	START="start.sh"
	LOG="output.log"

	# Kill any running process
	if ps ax | grep -v grep | grep $APP > /dev/null; then
		echo "SHUTDOWN RUNING APP..."
		killall -9 $APP
	fi

	cd $DEST
	echo "Fetching latest ZIP"
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/archive/$BRANCH.zip -O $BRANCH.zip
	echo "Extracting from $BRANCH.zip"
	unzip -oq $BRANCH.zip
	if [ -d $APP ]; then
		rm -r $APP
	fi
	mkdir -p $APP
	echo "Moving files from $REPO-$BRANCH/ to $APP"
	mv -f $REPO-$BRANCH/* $APP
	cd $APP
	chmod u+x $START
	echo "Starting $APP_TYPE in Mode=$APP_MODE"
	nohup sh $START --mode=$APP_MODE > $LOG 2>&1&
	echo "Done! $ sudo tail -f $DEST/$APP/$LOG for more details"
}

inform_loadbalancer () {
	retries=0
	if [[ "$1" == "REMOVE" ]]; then
		until curl -X DELETE http://$LOADBALANCER_HOST:$LOADBALANCER_API_PORT/$APP_TYPE/$APP_HOST:$APP_PORT; do
			echo "[$(timestamp)] Load Balancer is not running, retry... [$(( retries++ ))]"
			sleep 5
		done
	elif [[ "$1" == "ADD" ]]; then
		until curl -H "Content-Type: application/json" -X POST -d "{\"ip\":\"$APP_HOST:$APP_PORT\"}" http://$LOADBALANCER_HOST:$LOADBALANCER_API_PORT/$APP_TYPE; do
			echo "[$(timestamp)] Load Balancer is not up yet, retry... [$(( retries++ ))]"
			sleep 5
		done
	else
		echo >&2 "Error: Invalid argument";
		exit 1
	fi
}

echo "[$(timestamp)] ---- 1. Inform Load Balancer (Remove) ----"
inform_loadbalancer REMOVE

echo "[$(timestamp)] ---- 2. Run API (Master) Server ----"
run_master_server

echo "[$(timestamp)] ---- 3. Inform Load Balancer (Add) ----"
inform_loadbalancer ADD

echo "[$(timestamp)] ---- Completed ----"

exit 0
