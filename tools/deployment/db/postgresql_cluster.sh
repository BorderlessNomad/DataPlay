#!/bin/bash

# This is setup script for PostreSQL Cluster server.

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

restart_postgres () {
	systemctl restart postgresql-9.4
}

install_postgres () {
	rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
	yum update
	yum upgrade -y
	yum install -y postgresql94

	/usr/pgsql-9.4/bin/postgresql94-setup initdb

	systemctl start postgresql-9.4
	systemctl enable postgresql-9.4
}

setup_database () {
	cd # /var/lib/postgresql

	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"
	DB_NAME="dataplay"
	DB_VERSION="9.4"

	# Install PostgreSQL Adminpack
	psql --command "CREATE EXTENSION adminpack;"

	# Create a PostgreSQL user named 'playgen' with 'aDam3ntiUm' as the password and
	# then create a database 'dataplay' owned by the 'playgen' role.
	psql --command "CREATE USER $DB_USER WITH SUPERUSER PASSWORD '$DB_PASSWORD';" && \
	createdb -O $DB_USER $DB_NAME

	# Adjust PostgreSQL configuration so that remote connections to the database are possible.
	# From Private cluster & PlayGen dev IP
	echo "host    all             all             109.231.121.0/24        md5" >> /var/lib/pgsql/$DB_VERSION/data/pg_hba.conf
	echo "host    all             all             109.231.122.0/24        md5" >> /var/lib/pgsql/$DB_VERSION/data/pg_hba.conf
	echo "host    all             all             109.231.123.0/24        md5" >> /var/lib/pgsql/$DB_VERSION/data/pg_hba.conf
	echo "host    all             all             109.231.124.0/24        md5" >> /var/lib/pgsql/$DB_VERSION/data/pg_hba.conf
	echo "host    all             all             213.122.181.2/32        md5" >> /var/lib/pgsql/$DB_VERSION/data/pg_hba.conf

	# And add 'listen_addresses' to '/etc/postgresql/$DB_VERSION/main/postgresql.conf'
	echo "listen_addresses = '*'" >> /var/lib/pgsql/$DB_VERSION/data/postgresql.conf
	echo "port = 5432" >> /var/lib/pgsql/$DB_VERSION/data/postgresql.conf
}

install_pgpool () {
	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"

	yum install -y http://www.pgpool.net/yum/rpms/3.4/redhat/rhel-7-x86_64/pgpool-II-release-3.4-1.noarch.rpm
	yum install -y pgpool-II-94 pgpool-II-94-extensions

	cp /etc/pgpool-II-94/pcp.conf.sample /etc/pgpool-II-94/pcp.conf
	echo "$DB_USER:`pg_md5 $DB_PASSWORD`" >> /etc/pgpool-II-94/pcp.conf

	cp /etc/pgpool-II-94/pool_hba.conf.sample /etc/pgpool-II-94/pool_hba.conf
	echo "host    all         all         0.0.0.0/0             md5" >> /etc/pgpool-II-94/pool_hba.conf
}

update_iptables () {
	yum install -y firewalld
	systemctl start firewalld.service
	systemctl enable firewalld.service

	firewall-cmd --permanent --add-port=5432/tcp
	firewall-cmd --permanent --add-port=80/tcp

	firewall-cmd --reload
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install PostgresSQL ----"
install_postgres

echo "[$(timestamp)] ---- 3. Setup Database ----"
su postgres -c "$(typeset -f setup_database); setup_database" # Run function as user 'postgres'

echo "[$(timestamp)] ---- 4. Restart PostgresSQL ----"
restart_postgres

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
