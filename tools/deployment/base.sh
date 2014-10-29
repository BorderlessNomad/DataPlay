#!/bin/bash

# This is Base Image setup script

set -ex

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
	cat authorized_keys >> /home/ubuntu/.ssh/authorized_keys
}

update () {
	apt-get update
	apt-get -y upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo vim openssh-server gcc curl git mercurial bzr make binutils bison wget axel python-software-properties htop unzip
}

install_nodejs () {
	apt-add-repository -y ppa:chris-lea/node.js
	apt-get update
	apt-get install -y python g++ make nodejs
	npm install grunt-cli coffee-script bower -g
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
