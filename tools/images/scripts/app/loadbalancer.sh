#!/bin/bash

# This is setup script for Load Balancer.

set -ex

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
	npm cache clean && npm install -g coffee-script forever 2> /dev/null

	command -v haproxy 2>/dev/null || { echo >&2 "Error: Command 'haproxy' not found!"; exit 1; }

	command -v forever >/dev/null 2>&1 || { echo >&2 "Error: 'forever' is not installed!"; exit 1; }

	command -v coffee >/dev/null 2>&1 || { echo >&2 "Error: 'coffee' is not installed!"; exit 1; }

	mkdir -p /home/ubuntu && cd /home/ubuntu
	mkdir -p haproxy-api && cd haproxy-api

	wget -Nq https://raw.githubusercontent.com/playgenhub/DataPlay/develop/tools/images/scripts/app/haproxy-api/app.coffee && \
	wget -Nq https://raw.githubusercontent.com/playgenhub/DataPlay/develop/tools/images/scripts/app/haproxy-api/package.json && \
	wget -Nq https://raw.githubusercontent.com/playgenhub/DataPlay/develop/tools/images/scripts/app/haproxy-api/backend.json && \
	wget -Nq https://raw.githubusercontent.com/playgenhub/DataPlay/develop/tools/images/scripts/app/haproxy-api/haproxy.cfg.template

	npm install

	coffee -cb app.coffee > app.js

	forever start -l forever.log -o output.log -e errors.log app.js 2>/dev/null
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 1936 -j ACCEPT # HAProxy statistics
	iptables -A INPUT -p tcp --dport 1937 -j ACCEPT # HAProxy API

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

command -v node 2>/dev/null || { echo >&2 "Error: Command 'node' not found!"; exit 1; }

command -v npm 2>/dev/null || { echo >&2 "Error: Command 'npm' not found!"; exit 1; }

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
