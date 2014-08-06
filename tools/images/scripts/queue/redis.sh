#!/bin/bash

# This is setup script for Queue (RabbitMQ) + Non-persistent DV (Redis)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-redis-rabbitmq for Redis+RabbitMQ instance

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

update () {
	apt-get update
	apt-get -y upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo openssh-server screen gcc curl git make binutils bison wget python-software-properties htop zip
}

install_redis () {
	add-apt-repository -y ppa:rwky/redis
	apt-get update
	apt-get install -y redis-server
	wget -O /etc/init.d/redis-server "https://gist.githubusercontent.com/lsbardel/257298/raw/d48b84d89289df39eaddc53f1e9a918f776b3074/redis-server-for-init.d-startup"
	chmod 755 /etc/init.d/redis-server
	update-rc.d redis-server defaults
	service redis-server start
}

install_rabbitmq () {
	echo "deb http://www.rabbitmq.com/debian/ testing main" > /etc/apt/sources.list.d/rabbitmq.list
	wget http://www.rabbitmq.com/rabbitmq-signing-key-public.asc
	apt-key add rabbitmq-signing-key-public.asc
	apt-get update
	apt-get install -y rabbitmq-server
	/usr/sbin/rabbitmq-plugins enable rabbitmq_management
	echo "[{rabbit, [{loopback_users, []}]}]." > /etc/rabbitmq/rabbitmq.config
}

update_iptables () {
	# Monitoring ports 80, 8080, 4242, 4243, 4245
	iptables -A INPUT -p tcp --dport 80 -j ACCEPT
	iptables -A INPUT -p tcp --dport 443 -j ACCEPT
	iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4242 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4243 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4245 -j ACCEPT

	# Redis Ports
	iptables -A INPUT -p tcp --dport 6379 -j ACCEPT

	# RabbitMQ Ports
	iptables -A INPUT -p tcp --dport 5672 -j ACCEPT
	iptables -A INPUT -p tcp --dport 15672 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4369 -j ACCEPT

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

# As root
echo "---- Running as Root ----"

echo "1. ---- Setup Host ----"
setuphost
echo "2. ---- Update system ----"
update
echo "3. ---- Install essential packages ----"
install_essentials
echo "4. ---- Install Redis ----"
install_redis
echo "4. ---- Install RabbitMQ ----"
install_rabbitmq
echo "5. ---- Update IPTables rules ----"
update_iptables

echo "---- Completed ----"
