#!/bin/bash

# This is setup script for Load Balancer.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

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
	apt-get update > /dev/null
	apt-get install -y haproxy
}

setup_haproxy_api () {
	URL="https://raw.githubusercontent.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO/$BRANCH"

	npm cache clean && npm install -g coffee-script forever 2> /dev/null

	command -v haproxy >/dev/null 2>&1 || { echo >&2 "Error: Command 'haproxy' not found!"; exit 1; }

	command -v forever >/dev/null 2>&1 || { echo >&2 "Error: 'forever' is not installed!"; exit 1; }

	command -v coffee >/dev/null 2>&1 || { echo >&2 "Error: 'coffee' is not installed!"; exit 1; }

	mkdir -p /home/ubuntu && cd /home/ubuntu
	mkdir -p haproxy-api && cd haproxy-api

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/app.coffee && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/package.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/backend.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/loadbalancer/api/haproxy.cfg.template

	npm install

	coffee -cb app.coffee > app.js

	forever start -l forever.log -o output.log -e errors.log app.js >/dev/null 2>&1

	# curl -i -H "Accept: application/json" http://109.231.121.47:1937
	# curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X POST -d '{"ip":"109.231.121.61:80"}' http://109.231.121.47:1937
	# curl -i -H "Accept: application/json" -X DELETE http://109.231.121.47:1937/109.231.121.61:80
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 1936 -j ACCEPT # HAProxy statistics
	iptables -A INPUT -p tcp --dport 1937 -j ACCEPT # HAProxy API

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

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
