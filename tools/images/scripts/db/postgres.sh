#!/bin/bash

# This is setup script for PostreSQL Database server.

set -ex

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

install_postgres () {
	apt-get install -y postgresql postgresql-contrib postgresql-client libpq-dev
	apt-get autoclean
	# TODO setup users etc
}

setup_database () {
	echo "TODO"
	# TODO
	# pg_dump -v -cC -f dataplay.sql -h 10.0.0.2 -U playgen dataplay
	# pg_dump -v -cC -f dataplay.sql -h localhost -U playgen dataplay
	# gzip -vk dataplay.sql
	# scp dataplay.sql.gz ubuntu@109.231.121.12:/home/ubuntu
	# gunzip -vk dataplay.sql.gz
	# psql -h localhost -U playgen -d dataplay -f dataplay.sql
}

update_iptables () {
	iptables -A INPUT -p tcp --dport 5432 -j ACCEPT # PostgreSQL listener

	iptables-save
}

if [ "$(id -u)" != "0" ]; then
	echo "Error: This script must be run as root" 1>&2
	exit 1
fi

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install PostgresSQL ----"
install_postgres

echo "[$(timestamp)] ---- 3. Export Database ----"
setup_database

echo "[$(timestamp)] ---- 4. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
