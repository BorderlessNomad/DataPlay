#!/bin/bash

# This is setup script for Load Balancer.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

HOST=$(ifconfig eth0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}')
PORT="1938"

# REDIS_HOST="109.231.121.13"
REDIS_HOST=$(ss-get --timeout 360 redis.hostname)
REDIS_PORT="6379"

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

install_haproxy () {
	apt-add-repository -y ppa:vbernat/haproxy-1.5
	sed -i "s/utopic/trusty/g" /etc/apt/sources.list.d/vbernat-ubuntu-haproxy-1_5-utopic.list # ppa doesn't have support for 14.10 yet
	apt-get update
	apt-get install -y haproxy

	# Using single quotes to avoid bash $ variable expansion
	echo '# HAProxy' >> /etc/rsyslog.conf
	echo '$ModLoad imudp' >> /etc/rsyslog.conf
	echo '$UDPServerAddress 127.0.0.1' >> /etc/rsyslog.conf
	echo '$UDPServerRun 514' >> /etc/rsyslog.conf

	service rsyslog restart
	service haproxy restart
}

setup_haproxy_api () {
	URL="https://raw.githubusercontent.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO/$BRANCH"

	command -v haproxy >/dev/null 2>&1 || { echo >&2 "Error: Command 'haproxy' not found!"; exit 1; }

	command -v npm >/dev/null 2>&1 || { echo >&2 'Error: Command "npm" not found!'; exit 1; }

	command -v forever >/dev/null 2>&1 || { echo >&2 "Error: 'forever' is not installed!"; exit 1; }

	command -v coffee >/dev/null 2>&1 || { echo >&2 "Error: 'coffee-script' is not installed!"; exit 1; }

	mkdir -p /home/ubuntu && cd /home/ubuntu
	mkdir -p haproxy-api && cd haproxy-api

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/app.coffee && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/package.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/proxy.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/haproxy.cfg.template

	npm install

	coffee -cb app.coffee > app.js

	forever start -l forever.log -o output.log -e errors.log app.js >/dev/null 2>&1

	# Gamification:
	# curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X POST -d '{"ip":"109.231.121.55:80"}' http://109.231.121.84:1937/gamification
	# curl -i -H "Accept: application/json" -X DELETE http://109.231.121.84:1937/gamification/109.231.121.55:80
	#
	# Master:
	# curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X POST -d '{"ip":"109.231.121.94:3000"}' http://109.231.121.84:1937/master
	# curl -i -H "Accept: application/json" -X DELETE http://109.231.121.84:1937/master/109.231.121.94:3000
}

run_monitoring () {
	URL="https://github.com"
	USER="playgenhub"
	REPO="DataPlay-Monitoring"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO"
	DEST="/home/ubuntu/www"
	APP="dataplay-monitoring"

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
	echo "Starting $START"
	nohup sh $START > $LOG 2>&1&
	echo "Done! $ sudo tail -f $DEST/$APP/$LOG for more details"
}

export_variables () {
	echo "export DP_REDIS_HOST=$REDIS_HOST" >> /etc/profile.d/dataplay.sh
	echo "export DP_REDIS_PORT=$REDIS_PORT" >> /etc/profile.d/dataplay.sh
	echo "export DP_MONITORING_PORT=$PORT" >> /etc/profile.d/dataplay.sh

	. /etc/profile

	su - ubuntu -c ". /etc/profile"
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 1936 -j ACCEPT # HAProxy statistics
	iptables -A INPUT -p tcp --dport 1937 -j ACCEPT # HAProxy API
	iptables -A INPUT -p tcp --dport $PORT -j ACCEPT # API Health monitor

	iptables-save
}

command -v node >/dev/null 2>&1 || { echo >&2 "Error: Command 'node' not found!"; exit 1; }

command -v npm >/dev/null 2>&1 || { echo >&2 "Error: Command 'npm' not found!"; exit 1; }

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install HAProxy ----"
install_haproxy

echo "[$(timestamp)] ---- 3. Setup HAProxy API ----"
setup_haproxy_api

echo "[$(timestamp)] ---- 4. Export Variables ----"
export_variables

echo "[$(timestamp)] ---- 5. Run API Monitoring Probe ----"
run_monitoring

echo "[$(timestamp)] ---- 6. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
