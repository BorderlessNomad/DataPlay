#!/bin/bash

set -ex

# This is setup script for DB (PostreSQL server)
# 1. Install Ubuntu base image or dataplay-ubuntu-base (recommended)
# 2. Run this script as 'sudo'
#
# Note: Installing from pre-configured base image is highly recommended
#	e.g.
#		dataplay-posgresql for DB instance

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

install_postgres () {
	apt-get install -y postgresql postgresql-contrib postgresql-9.3-postgis-2.1 postgresql-client libpq-dev
	apt-get clean && sudo rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
}

sync_database () {
	# TODO
	# pg_dump -v -cC -f dataplay.sql -h 10.0.0.2 -U playgen dataplay
	# pg_dump -v -cC -f dataplay.sql -h localhost -U playgen dataplay
	# gzip -vk dataplay.sql
	# scp dataplay.sql.gz ubuntu@109.231.121.12:/home/ubuntu
	# gunzip -vk dataplay.sql.gz
	# psql -h localhost -U playgen -d dataplay -f dataplay.sql
}

update_iptables () {
	# Monitoring ports 80, 8080, 4242, 4243, 4245
	iptables -A INPUT -p tcp --dport 80 -j ACCEPT
	iptables -A INPUT -p tcp --dport 443 -j ACCEPT
	iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4242 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4243 -j ACCEPT
	iptables -A INPUT -p tcp --dport 4245 -j ACCEPT

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
echo "4. ---- Install PostgresSQL ----"
install_postgres
echo "5. ---- Update IPTables rules ----"
update_iptables

echo "---- Completed ----"

exit 0
