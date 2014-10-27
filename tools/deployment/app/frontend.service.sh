#!/bin/bash

# This is service script for Frontend instance.
# Used for updating app on server after deployment.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

DEST="/home/ubuntu/www"
APP="dataplay"
WWW="www-src"

APP_HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
APP_PORT="80"
APP_TYPE="gamification"

# . /etc/profile

# LOADBALANCER_HOST="109.231.121.26"
LOADBALANCER_HOST="$DP_LOADBALANCER_HOST"
LOADBALANCER_REQUEST_PORT="$DP_LOADBALANCER_REQUEST_PORT"
LOADBALANCER_API_PORT="$DP_LOADBALANCER_API_PORT"
DOMAIN="$DP_DOMAIN"
# DOMAIN="dataplay.playgen.com"
# DOMAIN="$LOADBALANCER_HOST:LOADBALANCER_REQUEST_PORT"

timestamp () {
	date +"%F %T,%3N"
}

download_app () {
	URL="https://github.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO"

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
}

init_frontend () {
	sed -i "s/localhost:3000/$DOMAIN/g" $DEST/$APP/$WWW/dist/scripts/*.js
}

reload_nginx () {
	service nginx reload
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

echo "[$(timestamp)] ---- 2. Download Application ----"
download_app

echo "[$(timestamp)] ---- 3. Init Frotnend ----"
init_frontend

echo "[$(timestamp)] ---- 4. Reload Nginx ----"
init_frontend

echo "[$(timestamp)] ---- 5. Inform Load Balancer (Add) ----"
inform_loadbalancer ADD

echo "[$(timestamp)] ---- Completed ----"

exit 0
