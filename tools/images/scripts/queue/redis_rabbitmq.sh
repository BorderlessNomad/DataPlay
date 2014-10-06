#!/bin/bash

# This is setup script for Queuing (RabbitMQ) + Non-persistent Key-Value store (Redis) server.

set -ex

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

	add-apt-repository -y ppa:rwky/redis
	apt-get update
	apt-get install -y redis-server

	service redis-server restart

	npm install -g redis-commander
	redis-commander > redis-commander.log &
}

install_rabbitmq () {
	mkdir -p /home/ubuntu && cd /home/ubuntu

	curl http://www.rabbitmq.com/rabbitmq-signing-key-public.asc | sudo apt-key add -
	echo "deb http://www.rabbitmq.com/debian/ testing main" > /etc/apt/sources.list.d/rabbitmq.list
	apt-get update
	apt-get install -y rabbitmq-server

	rabbitmqctl add_user playgen aDam3ntiUm && \
	rabbitmqctl set_permissions -p / playgen ".*" ".*" ".*" && \
	rabbitmqctl set_user_tags playgen administrator && \
	rabbitmqctl delete_user guest

	service rabbitmq-server restart
}

enable_rabbitmqadmin () {
	rabbitmq-plugins enable rabbitmq_management
	echo "[{rabbit, [{loopback_users, []}]}]." > /etc/rabbitmq/rabbitmq.config

	service rabbitmq-server restart > rabbitmq-service.log & # Start RabbitMQ in background
	echo "Waiting for RabbitMQ to restart..."
	sleep 1
	while ! grep -m1 '...done.' < rabbitmq-service.log ; do
		sleep 1
	done
	echo "RabbitMQ is UP!"

	while [ 1 ]; do
		wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 http://localhost:15672/cli/rabbitmqadmin
		if [ $? = 0 ]; then break; fi; # check return value, break if successful (0)
		sleep 1s;
	done;
	chmod +x rabbitmqadmin
	mv rabbitmqadmin /usr/local/sbin

	service rabbitmq-server restart
}

update_iptables () {
	# Redis Ports
	iptables -A INPUT -p tcp --dport 6379 -j ACCEPT # Socket connections
	iptables -A INPUT -p tcp --dport 8081 -j ACCEPT # Redis Commander

	# RabbitMQ Ports
	iptables -A INPUT -p tcp --dport 4369 -j ACCEPT # Erlang Port Mapper Daemon (epmd)
	iptables -A INPUT -p tcp --dport 5672 -j ACCEPT # Message Queue
	iptables -A INPUT -p tcp --dport 15672 -j ACCEPT # RabbitMQ Management console
	iptables -A INPUT -p tcp --dport 35197 -j ACCEPT # Cluster communication

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install Redis ----"
install_redis

echo "[$(timestamp)] ---- 3. Install RabbitMQ ----"
install_rabbitmq

echo "[$(timestamp)] ---- 4. Enable RabbitMQ Admin ----"
# enable_rabbitmqadmin # Must be installed manually!!!

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
