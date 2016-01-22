#!/bin/bash

# This is Base Image setup script

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
	HOSTLOCAL="127.0.0.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

setup_ssh_keys () {
	URL="https://raw.githubusercontent.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO/$BRANCH"

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/authorized_keys
	cat authorized_keys >> /home/centos/.ssh/authorized_keys
}

update () {
	yum update -y
	yum upgrade -y
}

install_essentials () {
	yum install -y htop wget rsyslog gcc-c++ make firewalld

	systemctl start firewalld.service
	systemctl enable firewalld.service
}

install_nodejs () {
	curl -sL https://rpm.nodesource.com/setup | bash -
	yum install -y nodejs
	npm install -g grunt-cli coffee-script bower forever
}

update_firewall () {
	firewall-cmd --permanent --add-port=80/tcp
	firewall-cmd --permanent --add-port=443/tcp
	firewall-cmd --permanent --add-port=8080/tcp

	firewall-cmd --reload
}

download_scripts () {
	mkdir -p scripts && cd scripts
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/db/cassandra.sh && chmod +x cassandra.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/db/pgpool.sh && chmod +x pgpool.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/db/redis.sh && chmod +x redis.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/db/postgresql.sh && chmod +x postgresql.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/app/master.sh && chmod +x master.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/app/frontend.sh && chmod +x frontend.sh
	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 $SOURCE/tools/deployment/loadbalancer/haproxy.sh && chmod +x haproxy.sh
}

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

echo "[$(timestamp)] ---- 2. Update system ----"
update

echo "[$(timestamp)] ---- 3. Install essential packages ----"
install_essentials

echo "[$(timestamp)] ---- 4. Setup SSH Access Keys ----"
setup_ssh_keys

echo "[$(timestamp)] ---- 5. Install Node.js ----"
install_nodejs

echo "[$(timestamp)] ---- 6. Update Firewall rules ----"
update_firewall

echo "[$(timestamp)] ---- 7. Download Application Scripts ----"
download_scripts

echo "[$(timestamp)] ---- Completed ----"

exit 0
