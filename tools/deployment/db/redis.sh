#!/bin/bash

# This is setup script for Redis server.

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

install_redis () {
	mkdir -p /home/ubuntu && cd /home/ubuntu

	apt-get update
	apt-get install -y redis-server

	cp /etc/redis/redis.conf /etc/redis/redis.conf.backup
	sed -i "s/bind 127.0.0.1/bind 0.0.0.0/" /etc/redis/redis.conf # Allow external connections

	/etc/init.d/redis-server restart
}

install_redis_admin () {
	npm install -g redis-commander
	nohup redis-commander > redis-commander.log 2>&1&
}

update_iptables () {
	# Redis Ports
	iptables -A INPUT -p tcp --dport 6379 -j ACCEPT # Socket connections
	iptables -A INPUT -p tcp --dport 8081 -j ACCEPT # Redis Commander

	iptables-save
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install Redis ----"
install_redis

echo "[$(timestamp)] ---- 3. Install Redis Admin ----"
install_redis_admin

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
