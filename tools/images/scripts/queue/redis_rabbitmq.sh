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
	add-apt-repository -y ppa:rwky/redis
	apt-get update
	apt-get install -y redis-server python python-dev python-pip python-virtualenv

	service redis-server restart

	npm install -g redis-commander
	redis-commander > redis-commander.log &
}

install_rabbitmq () {
	curl http://www.rabbitmq.com/rabbitmq-signing-key-public.asc | sudo apt-key add -
	echo "deb http://www.rabbitmq.com/debian/ testing main" > /etc/apt/sources.list.d/rabbitmq.list
	apt-get update
	apt-get install -y rabbitmq-server

	rabbitmq-plugins enable rabbitmq_management
	echo "[{rabbit, [{loopback_users, []}]}]." > /etc/rabbitmq/rabbitmq.config

	service rabbitmq-server restart

	wget http://localhost:15672/cli/rabbitmqadmin
	chmod +x rabbitmqadmin
	mv rabbitmqadmin /usr/local/sbin

	rabbitmqctl add_user playgen aDam3ntiUm
	rabbitmqctl set_permissions -p / playgen ".*" ".*" ".*"
	rabbitmqctl set_user_tags playgen administrator
	rabbitmqctl delete_user guest

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

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
