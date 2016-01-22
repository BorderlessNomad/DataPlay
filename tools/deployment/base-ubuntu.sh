#!/bin/bash

# This is Base Image setup script for Ubuntu 15.04

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

	wget --retry-connrefused --waitretry=1 --read-timeout=20 --timeout=15 -t 0 -N $SOURCE/tools/deployment/authorized_keys
	cat authorized_keys >> /home/ubuntu/.ssh/authorized_keys
}

update () {
	apt-get update
	sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" upgrade
}

install_essentials () {
	apt-get install -y build-essential sudo ntpdate vim openssh-server gcc curl git mercurial bzr make binutils bison wget axel python-software-properties htop unzip
}

install_nodejs () {
	curl -sL https://raw.githubusercontent.com/playgenhub/DataPlay/master/tools/deployment/nodejs.sh | bash -
	apt-get install -y python g++ make nodejs
	npm install -g grunt-cli coffee-script bower forever
}

update_firewall () {
	iptables -A INPUT -p tcp --dport 80 -j ACCEPT
	iptables -A INPUT -p tcp --dport 443 -j ACCEPT
	iptables -A INPUT -p tcp --dport 8080 -j ACCEPT

	iptables-save
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
