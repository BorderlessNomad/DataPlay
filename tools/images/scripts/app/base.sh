#!/bin/bash

# This is Base Image setup script

set -ex

timestamp () {
	date +"%F %T,%3N"
}

setuphost () {
	HOSTNAME=$(hostname)
	HOSTLOCAL="127.0.1.1"
	echo "$HOSTLOCAL $HOSTNAME" >> /etc/hosts
}

update () {
	apt-get update > /dev/null
	apt-get -y upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo vim openssh-server gcc curl git make binutils bison wget python-software-properties htop unzip
}

install_nodejs () {
	apt-add-repository -y ppa:chris-lea/node.js
	apt-get update > /dev/null
	apt-get install -y python g++ make nodejs
}

update_iptables () {
	# Monitoring ports 80, 8080, 4242, 4243, 4245 for JCatascopia
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

echo "[$(timestamp)] ---- 1. Setup Host ----"
setuphost

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
