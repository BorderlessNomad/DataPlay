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

install_pgpool () {
	DB_USER="playgen"
	DB_PASSWORD="aDam3ntiUm"

	yum install -y http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
	yum install -y http://www.pgpool.net/yum/rpms/3.4/redhat/rhel-7-x86_64/pgpool-II-release-3.4-1.noarch.rpm

	yum update -y

	yum install -y wget rsyslog pgpool-II-94 pgpool-II-94-extensions postgresql94

	sed -i 's/^#$ModLoad imudp/$ModLoad imudp/g' /etc/rsyslog.conf
	sed -i 's/^#$UDPServerRun 514/$UDPServerRun 514/g' /etc/rsyslog.conf
	sed -i 's/^#$ModLoad imtcp/$ModLoad imtcp/g' /etc/rsyslog.conf
	sed -i 's/^#$InputTCPServerRun 514/$InputTCPServerRun 514/g' /etc/rsyslog.conf
	echo "local0.*                                                /var/log/pgpool.log" >> /etc/rsyslog.conf

	systemctl restart rsyslog.service
	systemctl enable rsyslog.service

	/usr/pgsql-9.4/bin/postgresql94-setup initdb

	systemctl restart postgresql-9.4
	systemctl enable postgresql-9.4

	cp /etc/pgpool-II-94/pcp.conf.sample /etc/pgpool-II-94/pcp.conf
	echo "$DB_USER:`pg_md5 $DB_PASSWORD`" >> /etc/pgpool-II-94/pcp.conf

	cp /etc/pgpool-II-94/pool_hba.conf.sample /etc/pgpool-II-94/pool_hba.conf
	echo "host    all         all         0.0.0.0/0             md5" >> /etc/pgpool-II-94/pool_hba.conf

	pg_md5 -m -u $DB_USER $DB_PASSWORD # Generate pool_passwd

	systemctl restart pgpool-II-94
	systemctl enable pgpool-II-94
}

install_nodejs () {
	curl -sL https://rpm.nodesource.com/setup | bash -
	yum install -y nodejs
	# yum install gcc-c++ make
}

setup_pgpool_api () {
	URL="https://raw.githubusercontent.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO/$BRANCH"

	npm cache clean
	npm install -g coffee-script forever

	command -v pgpool >/dev/null 2>&1 || { echo >&2 'Error: Command "pgpool" not found!'; exit 1; }

	command -v forever >/dev/null 2>&1 || { echo >&2 'Error: "forever" is not installed!'; exit 1; }

	command -v coffee >/dev/null 2>&1 || { echo >&2 'Error: "coffee-script" is not installed!'; exit 1; }

	cd /root && mkdir -p pgpool-api && cd pgpool-api

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/app.coffee && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/package.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/cluster.json && \
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/db/api/pgpool.conf.template

	npm install

	coffee -cb app.coffee > app.js

	forever -a start -l forever.log -o output.log -e errors.log app.js >/dev/null 2>&1

	###
	# curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X POST -d '{"ip":"109.231.124.136"}' http://109.231.124.122:1937
	# curl -i -H "Accept: application/json" -X DELETE http://109.231.124.122:1937/109.231.124.136
	###
}

update_iptables () {
	yum install -y firewalld
	systemctl start firewalld.service
	systemctl enable firewalld.service

	firewall-cmd --permanent --add-port=1937/tcp # pgpool-API
	firewall-cmd --permanent --add-port=9999/tcp # pgpool

	# JCatascopia
	firewall-cmd --permanent --add-port=80/tcp
	firewall-cmd --permanent --add-port=8080/tcp
	firewall-cmd --permanent --add-port=4242/tcp
	firewall-cmd --permanent --add-port=4243/tcp
	firewall-cmd --permanent --add-port=4245/tcp

	firewall-cmd --reload
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Install pgpool-II ----"
install_pgpool

echo "[$(timestamp)] ---- 3. Install Node.js ----"
install_nodejs

echo "[$(timestamp)] ---- 4. Setup pgpool API ----"
setup_pgpool_api

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
