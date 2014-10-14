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
	apt-get install -y axel postgresql postgresql-contrib postgresql-client libpq-dev
	apt-get autoclean
	service postgresql restart
}

setup_database () {
	cd
	# Create a PostgreSQL user named 'playgen' with 'aDam3ntiUm' as the password and
	# then create a database 'dataplay' owned by the 'playgen' role.
	psql --command "CREATE USER playgen WITH SUPERUSER PASSWORD 'aDam3ntiUm';" && \
	createdb -O playgen dataplay

	# Adjust PostgreSQL configuration so that remote connections to the database are possible.
	echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf

	# And add 'listen_addresses' to '/etc/postgresql/9.3/main/postgresql.conf'
	echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf
}

import_data () {
	echo "localhost:5432:dataplay:playgen:aDam3ntiUm" > .pgpass && chmod 0600 .pgpass
	YESTERDAY=$(date +%Y-%m-%d) # Today
	BACKUP_HOST="108.61.197.87"
	BACKUP_PORT="8080"
	BACKUP_DIR="postgresql/$YESTERDAY-daily"
	BACKUP_FILE="dataplay.sql.gz"

	until axel -a "http://$BACKUP_HOST:$BACKUP_PORT/$BACKUP_DIR/$BACKUP_FILE"; do
		YESTERDAY=$(date +%Y-%m-%d --date="$YESTERDAY -1 days") # Decrement by 1 Day
		BACKUP_DIR="postgresql/$YESTERDAY-daily"
		echo "Latest backup not available, try fetching $YESTERDAY"
	done

	gunzip -vk dataplay.sql.gz
	nohup psql -h localhost -U playgen -d dataplay -f dataplay.sql > postgres-import.log &
	# on Dev
	# echo "10.0.0.2:5432:dataplay:playgen:aDam3ntiUm" > .pgpass
	# chmod 0600 .pgpass
	# pg_dump -v -cC -f dataplay.sql -h 10.0.0.2 -U playgen dataplay
	# gzip -vk dataplay.sql
	# scp dataplay.sql.gz ubuntu@109.231.121.12:/home/ubuntu
	#
	# on Server
	# gunzip -vk dataplay.sql.gz
	# psql -h localhost -U playgen -d dataplay -f dataplay.sql >> postgres-import.log
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

echo "[$(timestamp)] ---- 3. Setup Database ----"
su postgres -c "$(typeset -f setup_database); setup_database" # Run function as user 'postgres'

echo "[$(timestamp)] ---- 4. Import Data ----"
import_data

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
