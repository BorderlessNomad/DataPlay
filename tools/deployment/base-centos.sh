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

setup_ssh_keys () {
	URL="https://raw.githubusercontent.com"
	USER="playgenhub"
	REPO="DataPlay"
	BRANCH="master"
	SOURCE="$URL/$USER/$REPO/$BRANCH"

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/authorized_keys
	cat authorized_keys >> /home/centos/.ssh/authorized_keys
}

update () {
	yum update -y
	yum upgrade -y
}

install_essentials () {
	yum install -y htop wget rsyslog gcc-c++ make
}

install_nodejs () {
	curl -sL https://rpm.nodesource.com/setup | bash -
	yum install -y nodejs
	npm install -g grunt-cli coffee-script bower
}

update_iptables () {
	yum install -y firewalld
	systemctl start firewalld.service
	systemctl enable firewalld.service

	# JCatascopia
	firewall-cmd --permanent --add-port=80/tcp
	firewall-cmd --permanent --add-port=8080/tcp
	firewall-cmd --permanent --add-port=4242/tcp
	firewall-cmd --permanent --add-port=4243/tcp
	firewall-cmd --permanent --add-port=4245/tcp

	firewall-cmd --reload
}

echo "[$(timestamp)] ---- 1. Setup SSH Access Keys ----"
setup_ssh_keys

echo "[$(timestamp)] ---- 2. Update system ----"
update

echo "[$(timestamp)] ---- 3. Install essential packages ----"
install_essentials

echo "[$(timestamp)] ---- 4. Install Node.js ----"
install_nodejs

echo "[$(timestamp)] ---- 5. Update IPTables rules ----"
update_iptables

echo "[$(timestamp)] ---- Completed ----"

exit 0
