#!/bin/bash

# This is setup script for Load Balancer

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
	apt-get update
	apt-get -y upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo openssh-server gcc curl git make binutils bison wget python-software-properties htop zip
}

install_nodejs () {
	apt-add-repository -y ppa:chris-lea/node.js
	apt-get update
	apt-get install -y python g++ make nodejs
	npm install -g pm2 --unsafe-perm
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

# As root
echo "---- Running as Root ----"
timestamp

echo "1. ---- Setup Host ----"
timestamp
setuphost

echo "2. ---- Update system ----"
timestamp
update

echo "3. ---- Install essential packages ----"
timestamp
install_essentials

echo "4. ---- Install Node.js ----"
timestamp
install_nodejs

echo "5. ---- Update IPTables rules ----"
timestamp
update_iptables

echo "---- Completed ----"

timestamp

exit 0
