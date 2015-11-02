#!/bin/bash

# This is setup script for Redis server.

set -ex

if [ "$(id -u)" != "0" ]; then
	echo >&2 "Error: This script must be run as user 'root'";
	exit 1
fi

JCATASCOPIA_REPO="109.231.126.62"
JCATASCOPIA_DASHBOARD="109.231.122.112"

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

	service redis-server restart
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

#added to automate JCatascopiaAgent installation
setup_JCatascopiaAgent(){
	wget -q https://raw.githubusercontent.com/CELAR/celar-deployment/master/vm/jcatascopia-agent.sh

	bash ./jcatascopia-agent.sh > /tmp/JCata.txt 2>&1

	eval "sed -i 's/server_ip=.*/server_ip=$JCATASCOPIA_DASHBOARD/g' /usr/local/bin/JCatascopiaAgentDir/resources/agent.properties"

	#trying to solve issue with exists in restart and stop
	#screen -dmS JCata bash -c '/etc/init.d/JCatascopia-Agent stop  > /tmp/JCata.txt 2>&1'
	#sleep 2
	#screen -dmS JCata bash -c '/etc/init.d/JCatascopia-Agent start > /tmp/JCata.txt 2>&1'
	/etc/init.d/JCatascopia-Agent restart > /tmp/JCata.txt 2>&1

	rm ./jcatascopia-agent.sh
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install Redis ----"
install_redis

echo "[$(timestamp)] ---- 3. Install Redis Admin ----"
install_redis_admin

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- 5. Setting up JCatascopia Agent ----"
setup_JCatascopiaAgent

echo "[$(timestamp)] ---- Completed ----"

exit 0
