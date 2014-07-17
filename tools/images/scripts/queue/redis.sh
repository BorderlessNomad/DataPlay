#!/bin/bash

# This is setup script for Queue/Non-persistent (PostreSQL server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-redis-rabbitmq for Redis+RabbitMQ instance

HOSTNAME=$(hostname)
HOSTLOCAL="127.0.1.1"
echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts

apt-get update
apt-get -y dist-upgrade
apt-get install -y build-essential sudo openssh-server screen gcc curl git make binutils bison wget python-software-properties xclip

# Redis
add-apt-repository -y ppa:rwky/redis
apt-get update
apt-get install -y redis-server
wget -O /etc/init.d/redis-server "https://gist.githubusercontent.com/lsbardel/257298/raw/d48b84d89289df39eaddc53f1e9a918f776b3074/redis-server-for-init.d-startup"
chmod 755 /etc/init.d/redis-server
update-rc.d redis-server defaults
service redis-server start
